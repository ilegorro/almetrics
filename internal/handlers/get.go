package handlers

import (
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
	if err = tmpl.Execute(w, hctx.strg.GetMetrics()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Type", "charset=utf-8")
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
	w.Write([]byte(value))
}
