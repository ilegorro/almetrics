package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/handlers"
)

func MetricsRouter(hctx *handlers.HandlerContext) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{mType}/{mName}/{mValue}", WithLogging(hctx.UpdateHandler))
	r.Get("/value/{mType}/{mName}", WithLogging(hctx.GetValueHandler))
	r.Get("/", WithLogging(hctx.GetRootHandler))
	return r
}
