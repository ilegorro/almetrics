package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/common"
)

func (hctx *HandlerContext) GetRootHandler(w http.ResponseWriter, r *http.Request) {
	respHTML := `
		<html>
			<head>
				<title>Метрики</title>
			</head>
			<body>
				<ul>
					{{range $k, $v := .}}
						<li>{{$k}}: {{$v}}</li>
					{{end}}
				</ul>
			</body>
		</html> `
	tmpl, err := template.New("metrics").Parse(respHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err = tmpl.Execute(w, hctx.strg.GetMetrics()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hctx *HandlerContext) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	var value string
	switch mType {
	case common.MetricGauge:
		v, ok := hctx.strg.GetGauge(mName)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		value = fmt.Sprintf("%v", v)
	case common.MetricCounter:
		v, ok := hctx.strg.GetCounter(mName)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		value = fmt.Sprintf("%v", v)
	default:
		http.Error(w, "Incorrect type", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (hctx *HandlerContext) GetValueJSONHandler(w http.ResponseWriter, r *http.Request) {

	var data common.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(w, "Error parsing body", http.StatusInternalServerError)
		return
	}

	var respData common.Metrics
	switch data.MType {
	case common.MetricGauge:
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
