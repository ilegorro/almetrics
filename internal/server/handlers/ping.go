package handlers

import (
	"net/http"

	"github.com/ilegorro/almetrics/internal/server"
	"github.com/ilegorro/almetrics/internal/server/adapters/db"
)

func PingDBHandler(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dbAdapter, err := db.New(r.Context(), app.Options.DBDSN)
		if err != nil {
			http.Error(w, "Error connecting DB", http.StatusInternalServerError)
			return
		}
		err = dbAdapter.Pool.Ping(r.Context())
		if err != nil {
			http.Error(w, "Error ping DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
