package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ilegorro/almetrics/internal/server/config"
	"github.com/ilegorro/almetrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler(t *testing.T) {

	type reqParams struct {
		mType  string
		mName  string
		mValue string
	}
	tests := []struct {
		name   string
		params reqParams
		want   int
	}{
		{
			name: "correct gauge",
			params: reqParams{
				mType:  "gauge",
				mName:  "metric",
				mValue: "100",
			},
			want: http.StatusOK,
		},
		{
			name: "correct counter",
			params: reqParams{
				mType:  "counter",
				mName:  "metric",
				mValue: "100",
			},
			want: http.StatusOK,
		},
		{
			name: "incorrect gauge value",
			params: reqParams{
				mType:  "gauge",
				mName:  "metric",
				mValue: "wrong",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "incorrect counter value",
			params: reqParams{
				mType:  "counter",
				mName:  "metric",
				mValue: "wrong",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "incorrect type",
			params: reqParams{
				mType:  "foo",
				mName:  "metric",
				mValue: "100",
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			r := httptest.NewRequest(http.MethodPost, "/update", nil)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("mType", tt.params.mType)
			rctx.URLParams.Add("mName", tt.params.mName)
			rctx.URLParams.Add("mValue", tt.params.mValue)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			strg := storage.NewMemStorage()
			op := config.EmptyOptions()
			app := NewApp(strg, op)
			h := http.HandlerFunc(app.UpdateHandler)
			h(w, r)
			cancel()

			result := w.Result()

			assert.Equal(t, tt.want, result.StatusCode)

			result.Body.Close()
		})
	}
}

func TestUpdateJSONHandler(t *testing.T) {

	type reqParams struct {
		mType  string
		mName  string
		mvalue string
	}
	tests := []struct {
		name     string
		bodyJSON string
		want     int
	}{
		{
			name:     "correct gauge",
			bodyJSON: `{"id": "metric", "type": "gauge", "value": 100}`,
			want:     http.StatusOK,
		},
		{
			name:     "correct counter",
			bodyJSON: `{"id": "metric", "type": "counter", "delta": 100}`,
			want:     http.StatusOK,
		},
		{
			name:     "incorrect gauge value",
			bodyJSON: `{"id": "metric", "type": "gauge", "value": "wrong"}`,
			want:     http.StatusBadRequest,
		},
		{
			name:     "incorrect counter value",
			bodyJSON: `{"id": "metric", "type": "counter", "delta": "wrong"}`,
			want:     http.StatusBadRequest,
		},
		{
			name:     "incorrect type",
			bodyJSON: `{"id": "metric", "type": "foo", "value": 100}`,
			want:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			r := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer([]byte(tt.bodyJSON)))
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			strg := storage.NewMemStorage()
			op := config.EmptyOptions()
			app := NewApp(strg, op)
			h := http.HandlerFunc(app.UpdateJSONHandler)
			h(w, r)
			cancel()
			result := w.Result()

			assert.Equal(t, tt.want, result.StatusCode)
			result.Body.Close()
		})
	}
}
