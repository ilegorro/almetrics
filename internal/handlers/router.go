package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/middleware"
)

func MetricsRouter(hctx *HandlerContext) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.WithCompression)

	r.Post("/update/{mType}/{mName}/{mValue}", hctx.UpdateHandler)
	r.Post("/update/", hctx.UpdateJSONHandler)
	r.Get("/value/{mType}/{mName}", hctx.GetValueHandler)
	r.Post("/value/", hctx.GetValueJSONHandler)
	r.Get("/", hctx.GetRootHandler)
	return r
}
