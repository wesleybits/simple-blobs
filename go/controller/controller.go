package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/wesleybits/go-simple-globs/hooks"
	m "github.com/wesleybits/go-simple-globs/model"
)

type Controller interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id string)
	Put(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

type Item struct {
	Id    string             `json:"id"`
	Hooks []hooks.SimpleHook `json:"hooks"`
	Data  any                `json:"data"`
}

type ItemController struct {
	model m.ModelSpec[Item]
	wg *sync.WaitGroup
}

func NewBasicController(wg *sync.WaitGroup, model m.ModelSpec[Item]) *ItemController {
	return &ItemController{
		model: model,
		wg: wg,
	}
}

func notfound(w http.ResponseWriter) {
	w.WriteHeader(404)
}

func (c *ItemController) GetAll(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	data := c.model.GetAll()
	items := make([]Item, len(data))
	for i, d := range data {
		items[i] = d.Data
		items[i].Id = d.Id
	}
	enc.Encode(items)
	w.WriteHeader(200)
}

func (c *ItemController) Post(w http.ResponseWriter, r *http.Request) {
	var item Item
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	dec.Decode(&item)
	id := c.model.Create(item)
	item.Id = id
	enc.Encode(item)
	w.WriteHeader(200)
	go func() {
		c.wg.Add(1)
		hooks.CallAll(item.Hooks, item, hooks.Create)
		c.wg.Done()
	}()
}

func (c *ItemController) Get(w http.ResponseWriter, r *http.Request, id string) {
	item := c.model.Get(id)
	if item == nil {
		notfound(w)
		return
	}
	enc := json.NewEncoder(w)

	item.Data.Id = id
	enc.Encode(item.Data)
	w.WriteHeader(200)
}

func (c *ItemController) Put(w http.ResponseWriter, r *http.Request, id string) {
	if old := c.model.Get(id); old == nil {
		notfound(w)
		return
	}

	var item Item
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	dec.Decode(&item)
	c.model.Update(id, item)
	item.Id = id
	enc.Encode(item)
	w.WriteHeader(200)
	go func() {
		c.wg.Add(1)
		hooks.CallAll(item.Hooks, item, hooks.Update)
		c.wg.Done()
	}()
}

func (c *ItemController) Delete(w http.ResponseWriter, r *http.Request, id string) {
	item := c.model.Get(id)
	if item == nil {
		notfound(w)
		return
	}
	item.Data.Id = item.Id
	c.model.Delete(id)
	w.WriteHeader(200)
	go func() {
		c.wg.Add(1)
		hooks.CallAll(item.Data.Hooks, item.Data, hooks.Delete)
		c.wg.Done()
	}()
}
