package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/storage"
)

func (ctx *HandlerContext) GetRootHandler(w http.ResponseWriter, r *http.Request) {
	respHTML := `<html>
		<head>
			<title>Метрики</title>
		</head>
		<body>
			<ul>`
	for k, v := range ctx.strg.GetMetrics() {
		respHTML += fmt.Sprintf("<li>%v: %v</li>\n", k, v)
	}
	respHTML += `</ul>
		</body>
		</html> `
	w.Write([]byte(respHTML))
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Type", "charset=utf-8")
}

func (ctx *HandlerContext) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	var value string
	switch mType {
	case storage.MetricGauge:
		v, ok := ctx.strg.GetGauge(mName)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		value = fmt.Sprintf("%v", v)
	case storage.MetricCounter:
		v, ok := ctx.strg.GetCounter(mName)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		value = fmt.Sprintf("%v", v)
	default:
		http.Error(w, "Incorrect type", http.StatusBadRequest)
		return
	}
	w.Write([]byte(value))
	w.WriteHeader(http.StatusOK)
}
