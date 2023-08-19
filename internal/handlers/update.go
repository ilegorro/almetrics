package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/common"
)

func (hctx *HandlerContext) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mValue := chi.URLParam(r, "mValue")
	switch mType {
	case common.MetricGauge:
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "Incorrect value", http.StatusBadRequest)
			return
		}
		hctx.strg.AddGauge(mName, common.Gauge(val))
	case common.MetricCounter:
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, "Incorrect value", http.StatusBadRequest)
			return
		}
		hctx.strg.AddCounter(mName, common.Counter(val))
	default:
		http.Error(w, "Incorrect type", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hctx *HandlerContext) UpdateJSONHandler(w http.ResponseWriter, r *http.Request) {
	var data common.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	var respData common.Metrics
	switch data.MType {
	case common.MetricGauge:
		hctx.strg.AddGauge(data.ID, common.Gauge(*data.Value))
		v, ok := hctx.strg.GetGauge(data.ID)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		respData = common.Metrics{
			ID:    data.ID,
			MType: common.MetricGauge,
			Value: (*float64)(&v),
		}
	case common.MetricCounter:
		hctx.strg.AddCounter(data.ID, common.Counter(*data.Delta))
		v, ok := hctx.strg.GetCounter(data.ID)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		respData = common.Metrics{
			ID:    data.ID,
			MType: common.MetricCounter,
			Delta: (*int64)(&v),
		}
	default:
		http.Error(w, "Incorrect type", http.StatusBadRequest)
		return
	}
	respJSON, err := json.MarshalIndent(respData, "", "  ")
	if err != nil {
		http.Error(w, "Error writing body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respJSON))
}
