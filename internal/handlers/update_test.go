package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilegorro/almetrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func Test_updateHandlerContext_UpdateHandler(t *testing.T) {

	tests := []struct {
		name    string
		request string
		want    int
	}{
		{
			name:    "correct gauge",
			request: "/update/gauge/metric/100",
			want:    http.StatusOK,
		},
		{
			name:    "correct counter",
			request: "/update/counter/metric/100",
			want:    http.StatusOK,
		},
		{
			name:    "incorrect path",
			request: "/update/100",
			want:    http.StatusNotFound,
		},
		{
			name:    "incorrect gauge value",
			request: "/update/gauge/metric/wrong",
			want:    http.StatusBadRequest,
		},
		{
			name:    "incorrect counter value",
			request: "/update/counter/metric/wrong",
			want:    http.StatusBadRequest,
		},
		{
			name:    "incorrect type",
			request: "/update/foo/metric/wrong",
			want:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			strg := storage.NewMemStorage()
			hctx := NewUpdateHandlerContext(&strg)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(hctx.UpdateHandler)
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want, result.StatusCode)

			result.Body.Close()
		})
	}
}
