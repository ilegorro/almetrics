package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/handlers"
)

func MetricsRouter(hctx *handlers.HandlerContext) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{mType}/{mName}/{mValue}", WithLogging(hctx.UpdateHandler))
	r.Post("/update/", WithLogging(hctx.UpdateJSONHandler))
	r.Get("/value/{mType}/{mName}", WithLogging(hctx.GetValueHandler))
	r.Post("/value/", WithLogging(hctx.GetValueJSONHandler))
	r.Get("/", WithLogging(hctx.GetRootHandler))
	return r
}
