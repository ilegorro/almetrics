package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/handlers"
	"github.com/ilegorro/almetrics/internal/storage"
)

func metricsRouter(ctx handlers.HandlerContext) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{mType}/{mName}/{mValue}", ctx.UpdateHandler)
	r.Get("/value/{mType}/{mName}", ctx.GetValueHandler)
	r.Get("/", ctx.GetRootHandler)
	return r
}

func main() {
	strg := storage.NewMemStorage()
	hctx := handlers.NewHandlerContext(&strg)

	err := http.ListenAndServe(`:8080`, metricsRouter(*hctx))
	if err != nil {
		panic(err)
	}
}
