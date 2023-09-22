package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ilegorro/almetrics/internal/agent/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_Report(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "no error",
		},
	}
	ctx := context.Background()
	op := config.EmptyOptions()
	app := NewApp(op)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			app.PollMemStats(ctx)
			u, err := url.Parse(srv.URL)
			require.NoError(t, err)
			if err == nil {
				app.Options.Endpoint.Hostname = u.Hostname()
				app.Options.Endpoint.Port = u.Port()
			}
			err = app.Report(ctx)
			require.NoError(t, err)
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
	op := config.EmptyOptions()
	app := NewApp(op)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.PollMemStats(ctx)
			assert.NotEmpty(t, app.metrics)
		})
	}
}
