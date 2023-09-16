package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/middleware"
	"github.com/ilegorro/almetrics/internal/server"
)

func MetricsRouter(app *server.App) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.WithHash(app))
	r.Use(middleware.WithCompression)

	r.Post("/update/{mType}/{mName}/{mValue}", UpdateHandler(app))
	r.Post("/update/", UpdateJSONHandler(app))
	r.Post("/updates/", UpdatesHandler(app))
	r.Get("/value/{mType}/{mName}", GetValueHandler(app))
	r.Post("/value/", GetValueJSONHandler(app))
	r.Get("/ping", PingDBHandler(app))
	r.Get("/", GetRootHandler(app))
	return r
}
