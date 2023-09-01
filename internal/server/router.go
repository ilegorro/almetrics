package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/middleware"
)

func MetricsRouter(app *App) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.WithCompression)

	r.Post("/update/{mType}/{mName}/{mValue}", app.UpdateHandler)
	r.Post("/update/", app.UpdateJSONHandler)
	r.Get("/value/{mType}/{mName}", app.GetValueHandler)
	r.Post("/value/", app.GetValueJSONHandler)
	r.Get("/ping", app.PingDBHandler)
	r.Get("/", app.GetRootHandler)
	return r
}
