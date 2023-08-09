package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilegorro/almetrics/internal/handlers"
	"github.com/ilegorro/almetrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func TestMetricsRouter(t *testing.T) {
	strg := storage.NewMemStorage()
	hctx := handlers.NewHandlerContext(strg)
	ts := httptest.NewServer(MetricsRouter(hctx))
	defer ts.Close()

	var testTable = []struct {
		url    string
		method string
		want   int
	}{
		{"/update/gauge/foo/100", "POST", http.StatusOK},
		{"/value/gauge/foo", "GET", http.StatusOK},
		{"/", "GET", http.StatusOK},
		{"/update", "POST", http.StatusNotFound},
		{"/value", "GET", http.StatusNotFound},
		{"/update/gauge", "POST", http.StatusNotFound},
		{"/value/gauge", "GET", http.StatusNotFound},
		{"/update/gauge/foo", "POST", http.StatusNotFound},
	}
	for _, v := range testTable {
		resp := testRequest(t, ts, v.method, v.url)
		assert.Equal(t, v.want, resp.StatusCode)
		resp.Body.Close()
	}
}
