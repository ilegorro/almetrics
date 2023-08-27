package handlers

import (
	"bytes"
	"encoding/json"
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
	data := make(map[string]string, 0)
	for _, v := range hctx.strg.GetMetrics() {
		data[v.ID] = v.StringValue()
	}
	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hctx *HandlerContext) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	v, err := hctx.strg.GetMetric(mName, mType)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(v.StringValue()))
}

func (hctx *HandlerContext) GetValueJSONHandler(w http.ResponseWriter, r *http.Request) {

	var data common.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(w, "Error parsing body", http.StatusInternalServerError)
		return
	}

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
}
