package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/middleware"
)

func MetricsRouter(hctx *HandlerContext) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{mType}/{mName}/{mValue}", middleware.WithLogging(middleware.WithCompression(hctx.UpdateHandler)))
	r.Post("/update/", middleware.WithLogging(middleware.WithCompression(hctx.UpdateJSONHandler)))
	r.Get("/value/{mType}/{mName}", middleware.WithLogging(middleware.WithCompression(hctx.GetValueHandler)))
	r.Post("/value/", middleware.WithLogging(middleware.WithCompression(hctx.GetValueJSONHandler)))
	r.Get("/", middleware.WithLogging(middleware.WithCompression(hctx.GetRootHandler)))
	return r
}
