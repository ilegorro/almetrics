package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ilegorro/almetrics/internal/storage"
)

type updateHandlerContext struct {
	strg storage.Repository
}

func NewUpdateHandlerContext(strg storage.Repository) *updateHandlerContext {
	if strg == nil {
		panic("Storage is not defined")
	}
	return &updateHandlerContext{strg: strg}
}

func (ctx *updateHandlerContext) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		path := strings.TrimPrefix(r.URL.Path, "/")
		params := strings.Split(path, "/")
		if len(params) != 4 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		mType := params[1]
		mName := params[2]
		mValue := params[3]
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
	} else {
		http.Error(w, "Only POST method is allowed", http.StatusForbidden)
		return
	}
}
