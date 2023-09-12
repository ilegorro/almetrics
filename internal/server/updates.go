package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
)

func (app *App) UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	var data []common.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = app.strg.AddMetrics(ctx, data)
	cancel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if app.syncFileStorage {
		sop := filestorage.Options{StoragePath: app.options.Storage.Path}
		err := filestorage.SaveMetrics(app.strg, &sop)
		if err != nil {
			common.SugaredLogger().Errorf("Error saving metrics: %v", err)
		}
	}
}
