package config

import (
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/ilegorro/almetrics/internal/common"
)

type config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

type Endpoint struct {
	Hostname string
	Port     string
}

func (e *Endpoint) URL() string {
	return fmt.Sprintf("http://%v:%v/update/", e.Hostname, e.Port)
}

type Options struct {
	Endpoint       *Endpoint
	PollInterval   int
	ReportInterval int
	Key            string
	RateLimit      int
}

func EmptyOptions() *Options {
	return &Options{
		Endpoint: &Endpoint{},
	}
}

func ReadOptions() *Options {
	logger := common.SugaredLogger()
	op := EmptyOptions()
	op.Endpoint.Hostname = "localhost"
	op.Endpoint.Port = "8080"

	parseFlags(op)
	cfg, err := readEnv()
	if err != nil {
		logger.Errorf("error parsing environment variables: %+v", err)
	}

	if cfg.PollInterval != 0 {
		op.PollInterval = cfg.PollInterval
	}
	if cfg.ReportInterval != 0 {
		op.ReportInterval = cfg.ReportInterval
	}
	if cfg.Address != "" {
		op.Endpoint, err = getEndpoint(cfg.Address)
		if err != nil {
			logger.Fatalf("error parsing hostname and port: %+v", err)
		}
	}
	if cfg.Key != "" {
		op.Key = cfg.Key
	}
	if cfg.RateLimit != 0 {
		op.RateLimit = cfg.RateLimit
	}

	return op
}

func readEnv() (config, error) {
	var cfg config
	err := env.Parse(&cfg)

	return cfg, err
}

func parseFlags(op *Options) {
	flag.IntVar(&op.PollInterval, "p", 2, "poll interval")
	flag.IntVar(&op.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&op.RateLimit, "l", 3, "rate limit")
	flag.StringVar(&op.Key, "k", "", "hash key")
	flag.Func("a", "host and port (default localhost:8080)", func(flagValue string) error {
		v, err := getEndpoint(flagValue)
		if err != nil {
			return fmt.Errorf("parse a flag: %w", err)
		}
		op.Endpoint = v

		return nil
	})
	flag.Parse()
}

func getEndpoint(s string) (*Endpoint, error) {
	e := &Endpoint{}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("get endpoint: %w", err)
	}
	e.Hostname = u.Hostname()
	e.Port = u.Port()

	return e, nil
}
