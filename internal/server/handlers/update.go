package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
	"github.com/ilegorro/almetrics/internal/server"
)

func UpdateHandler(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		err := app.Strg.AddMetric(ctx, &metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		if app.SyncFileStorage {
			sop := filestorage.Options{StoragePath: app.Options.Storage.Path}
			err := filestorage.SaveMetrics(r.Context(), app.Strg, &sop)
			if err != nil {
				common.SugaredLogger().Errorf("Error saving metrics: %v", err)
			}
		}
	}
}

func UpdateJSONHandler(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		err = app.Strg.AddMetric(ctx, &data)
		cancel()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		ctx, cancel = context.WithTimeout(r.Context(), 10*time.Second)
		v, err := app.Strg.GetMetric(ctx, data.ID, data.MType)
		cancel()

		if err != nil {
			if errors.Is(err, common.ErrWrongMetricsID) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else if errors.Is(err, common.ErrWrongMetricsType) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, fmt.Sprintf("error getting metric: %v", err), http.StatusInternalServerError)
			}
			return
		}

		respJSON, err := json.Marshal(v)
		if err != nil {
			http.Error(w, "error writing body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respJSON))
		if app.SyncFileStorage {
			sop := filestorage.Options{StoragePath: app.Options.Storage.Path}
			err := filestorage.SaveMetrics(r.Context(), app.Strg, &sop)
			if err != nil {
				common.SugaredLogger().Errorf("Error saving metrics: %v", err)
			}
		}
	}
}
