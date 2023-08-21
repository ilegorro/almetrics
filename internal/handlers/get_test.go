package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/middleware"
	"github.com/ilegorro/almetrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValueJSONHandler(t *testing.T) {
	var testGauge float64 = 100
	var testCounter int64 = 100
	type resp struct {
		gauge   *float64
		counter *int64
	}
	tests := []struct {
		name    string
		metrics common.Metrics
		req     common.Metrics
		want    resp
	}{
		{
			name: "get gauge",
			metrics: common.Metrics{
				MType: common.MetricGauge,
				ID:    "test_gauge",
				Value: (*float64)(&testGauge),
			},
			req: common.Metrics{
				MType: common.MetricGauge,
				ID:    "test_gauge",
			},
			want: resp{
				gauge: (*float64)(&testGauge),
			},
		},
		{
			name: "get counter",
			metrics: common.Metrics{
				MType: common.MetricCounter,
				ID:    "test_counter",
				Delta: (*int64)(&testCounter),
			},
			req: common.Metrics{
				MType: common.MetricCounter,
				ID:    "test_counter",
			},
			want: resp{
				counter: (*int64)(&testCounter),
			},
		},
	}

	strg := storage.NewMemStorage()
	hctx := NewHandlerContext(strg)
	updateHandler := http.HandlerFunc(middleware.WithCompression(hctx.UpdateJSONHandler))
	valueHandler := http.HandlerFunc(middleware.WithCompression(hctx.GetValueJSONHandler))
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updateHandler)
	mux.HandleFunc("/value/", valueHandler)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// update metrics
			dataJSON, err := json.MarshalIndent(tt.metrics, "", "  ")
			require.NoError(t, err)
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err = zb.Write([]byte(dataJSON))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, srv.URL+"/update/", buf)
			r.RequestURI = ""
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)

			zr, err := gzip.NewReader(resp.Body)
			require.NoError(t, err)
			respJSON, err := io.ReadAll(zr)
			require.NoError(t, err)

			var data common.Metrics
			err = json.Unmarshal(respJSON, &data)
			require.NoError(t, err, respJSON)

			assert.Equal(t, data.Value, tt.want.gauge)
			assert.Equal(t, data.Delta, tt.want.counter)
			resp.Body.Close()

			//get value
			dataJSON, err = json.MarshalIndent(tt.req, "", "  ")
			require.NoError(t, err)
			buf = bytes.NewBuffer(nil)
			zb = gzip.NewWriter(buf)
			_, err = zb.Write([]byte(dataJSON))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)

			r = httptest.NewRequest(http.MethodPost, srv.URL+"/value/", buf)
			r.RequestURI = ""
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")
			resp, err = http.DefaultClient.Do(r)
			require.NoError(t, err)

			zr, err = gzip.NewReader(resp.Body)
			require.NoError(t, err)
			respJSON, err = io.ReadAll(zr)
			require.NoError(t, err)

			var dataValue common.Metrics
			err = json.Unmarshal(respJSON, &dataValue)
			require.NoError(t, err, respJSON)

			assert.Equal(t, dataValue.Value, tt.want.gauge)
			assert.Equal(t, dataValue.Delta, tt.want.counter)
			resp.Body.Close()
		})
	}
}
