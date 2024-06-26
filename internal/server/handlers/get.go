package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/server"
)

func GetRootHandler(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		metrics, err := app.Strg.GetMetrics(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make(map[string]string, 0)
		for _, v := range metrics {
			data[v.ID] = v.StringValue()
		}

		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err = tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetValueHandler(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		v, err := app.Strg.GetMetric(ctx, mName, mType)
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

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(v.StringValue()))
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting metric: %v", err), http.StatusInternalServerError)
		}
	}
}

func GetValueJSONHandler(app *server.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data common.Metrics
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
			http.Error(w, "Error parsing body", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		v, err := app.Strg.GetMetric(ctx, data.ID, data.MType)
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
			http.Error(w, "Error writing body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(respJSON))
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting metric: %v", err), http.StatusInternalServerError)
		}
	}
}
