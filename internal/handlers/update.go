package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/storage"
)

func (ctx *HandlerContext) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mValue := chi.URLParam(r, "mValue")
	switch mType {
	case storage.MetricGauge:
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "Incorrect value", http.StatusBadRequest)
			return
		}
		ctx.strg.AddGauge(mName, storage.Gauge(val))
	case storage.MetricCounter:
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, "Incorrect value", http.StatusBadRequest)
			return
		}
		ctx.strg.AddCounter(mName, storage.Counter(val))
	default:
		http.Error(w, "Incorrect type", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
