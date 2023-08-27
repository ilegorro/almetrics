package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
)

func (hctx *HandlerContext) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mValue := chi.URLParam(r, "mValue")
	switch mType {
	case common.MetricGauge:
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, common.ErrWrongMetricsValue.Error(), http.StatusBadRequest)
			return
		}
		hctx.strg.AddGauge(mName, common.Gauge(val))
	case common.MetricCounter:
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, common.ErrWrongMetricsValue.Error(), http.StatusBadRequest)
			return
		}
		hctx.strg.AddCounter(mName, common.Counter(val))
	default:
		http.Error(w, common.ErrWrongMetricsType.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	if hctx.syncPath != "" {
		sop := filestorage.Options{StoragePath: hctx.syncPath}
		err := filestorage.SaveMetrics(hctx.strg, &sop)
		if err != nil {
			common.SugaredLogger().Errorf("Error saving metrics: %v", err)
		}
	}
}

func (hctx *HandlerContext) UpdateJSONHandler(w http.ResponseWriter, r *http.Request) {
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

	hctx.strg.AddMetric(&data)
	v, err := hctx.strg.GetMetric(data.ID, data.MType)
	if err != nil {
		switch err {
		case common.ErrWrongMetricsName:
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
	if hctx.syncPath != "" {
		sop := filestorage.Options{StoragePath: hctx.syncPath}
		err := filestorage.SaveMetrics(hctx.strg, &sop)
		if err != nil {
			common.SugaredLogger().Errorf("Error saving metrics: %v", err)
		}
	}
}
