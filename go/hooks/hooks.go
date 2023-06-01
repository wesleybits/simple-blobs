package hooks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ChangeType string

const (
	Update ChangeType = "Update"
	Create ChangeType = "Create"
	Delete ChangeType = "Delete"
)

type WebHook interface {
	Client() (*http.Client, error)
	Request() (*http.Request, error)
}

type HookBody[T any] struct {
	Change ChangeType `json:"change"`
	Data   T          `json:"data"`
}

func Call[T any](hook WebHook, data T, changeType ChangeType) error {
	client, err := hook.Client()
	if err != nil {
		return err
	}

	request, err := hook.Request()
	if err != nil {
		return err
	}

	bodyStruct := HookBody[T]{Change: changeType, Data: data}
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	bodyWriter := io.NopCloser(bytes.NewBuffer(bodyBytes))
	request.Body = bodyWriter
	request.ContentLength = int64(len(bodyBytes))
	request.Header["Content-Type"] = []string{"application/json"}
	_, err = client.Do(request)
	return err
}

func CallAll[T any, U WebHook](hooks []U, data T, changeType ChangeType) error {
	errs := make([]error, 0, len(hooks))
	for _, hook := range hooks {
		err := Call(hook, data, changeType)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs[0]
}

type SimpleHook string

func (hook SimpleHook) Client() (*http.Client, error) {
	var client http.Client
	client = *http.DefaultClient
	return &client, nil
}

func (hook SimpleHook) Request() (*http.Request, error) {
	hookUrl, err := url.Parse(string(hook))
	if err != nil {
		return nil, err
	}

	return &http.Request{
		Method: "POST",
		URL:    hookUrl,
	}, nil
}
