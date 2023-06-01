package router

import (
	"net/http"
	"strings"

	"github.com/wesleybits/go-simple-globs/controller"
)

func nothandled(w http.ResponseWriter) {
	w.WriteHeader(405)
}

func AddPath(mux *http.ServeMux, prefix string, control controller.Controller) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		switch id := strings.TrimPrefix(r.URL.Path+"/", prefix); id {
		case "":
			switch r.Method {
			case "GET":
				control.GetAll(w, r)
			case "POST":
				control.Post(w, r)
			default:
				nothandled(w)
			}
		default:
			switch r.Method {
			case "GET":
				control.Get(w, r, id)
			case "PUT":
				control.Put(w, r, id)
			case "DELETE":
				control.Delete(w, r, id)
			default:
				nothandled(w)
			}
		}
	})
}
