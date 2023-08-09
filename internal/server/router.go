package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/handlers"
)

func MetricsRouter(hctx *handlers.HandlerContext) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{mType}/{mName}/{mValue}", hctx.UpdateHandler)
	r.Get("/value/{mType}/{mName}", hctx.GetValueHandler)
	r.Get("/", hctx.GetRootHandler)
	return r
}
