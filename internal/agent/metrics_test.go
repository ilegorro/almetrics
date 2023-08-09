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
	}{
		{
			name: "no error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			m := NewMetrics()
			err := m.Report(srv.URL)
			assert.NoError(t, err)
			srv.Close()
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
