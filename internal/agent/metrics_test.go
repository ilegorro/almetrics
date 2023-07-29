package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_Report(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "server response ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()
			m := NewMetrics()
			m.Report(srv.URL)
		})
	}
}

func TestMetrics_Poll(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "not empty result",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.Poll()
			assert.NotEmpty(t, m.counter)
			assert.NotEmpty(t, m.gauge)
		})
	}
}
