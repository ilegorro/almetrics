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
	app := NewApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			app.Poll()
			err := app.Report(srv.URL)
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
	app := NewApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.Poll()
			assert.NotEmpty(t, app.metrics)
		})
	}
}
