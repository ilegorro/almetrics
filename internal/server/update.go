package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
)

func (app *App) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mValue := chi.URLParam(r, "mValue")
	metrics := common.Metrics{ID: mName, MType: mType}
	switch mType {
	case common.MetricGauge:
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, common.ErrWrongMetricsValue.Error(), http.StatusBadRequest)
			return
		}
		metrics.Value = &val
	case common.MetricCounter:
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, common.ErrWrongMetricsValue.Error(), http.StatusBadRequest)
			return
		}
		metrics.Delta = &val
	default:
		http.Error(w, common.ErrWrongMetricsType.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := app.strg.AddMetric(ctx, &metrics)
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

func (app *App) UpdateJSONHandler(w http.ResponseWriter, r *http.Request) {
	var data common.Metrics
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
	err = app.strg.AddMetric(ctx, &data)
	cancel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	v, err := app.strg.GetMetric(ctx, data.ID, data.MType)
	cancel()
	if err != nil {
		switch err {
		case common.ErrWrongMetricsID:
			http.Error(w, err.Error(), http.StatusNotFound)
		case common.ErrWrongMetricsType:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "error getting metric", http.StatusInternalServerError)
		}
		return
	}

	respJSON, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Error writing body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respJSON))
	if app.syncFileStorage {
		sop := filestorage.Options{StoragePath: app.options.Storage.Path}
		err := filestorage.SaveMetrics(app.strg, &sop)
		if err != nil {
			common.SugaredLogger().Errorf("Error saving metrics: %v", err)
		}
	}
}
