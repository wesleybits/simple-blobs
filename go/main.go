package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/wesleybits/go-simple-globs/controller"
	"github.com/wesleybits/go-simple-globs/model"
	"github.com/wesleybits/go-simple-globs/repo"
	"github.com/wesleybits/go-simple-globs/router"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)

	r := repo.NewInMemoryRepo[controller.Item](ctx)
	m := model.NewModel[controller.Item](r)
	c := controller.NewBasicController(wg, m)
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	router.AddPath(mux, "/items", c)
	r.Start(wg)
	server.ListenAndServe()
	stop()
	wg.Wait()
}
