package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/ilegorro/almetrics/internal/common"
)

type config struct {
	Address         string `env:"ADDRESS"`
	StorageInterval int    `env:"STORE_INTERVAL"`
	StoragePath     string `env:"FILE_STORAGE_PATH"`
	StorageRestore  bool   `env:"RESTORE"`
	DBDSN           string `env:"DATABASE_DSN"`
}

type Options struct {
	Endpoint    *Endpoint
	Storage     *Storage
	DBDSN       string
	EndpointURL string
}

type Storage struct {
	Interval int
	Path     string
	Restore  bool
}

type Endpoint struct {
	Hostname string
	Port     string
}

func EmptyOptions() *Options {
	return &Options{
		Endpoint: &Endpoint{},
		Storage:  &Storage{},
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

	if cfg.StoragePath != "" {
		op.Storage.Path = cfg.StoragePath
	}
	if cfg.StorageInterval != 0 {
		op.Storage.Interval = cfg.StorageInterval
	}
	_, ok := os.LookupEnv("RESTORE")
	if ok {
		op.Storage.Restore = cfg.StorageRestore
	}

	if cfg.Address != "" {
		op.Endpoint, err = getEndpoint(cfg.Address)
		if err != nil {
			logger.Fatalf("error parsing hostname and port: %+v", err)
		}
	}
	if cfg.DBDSN != "" {
		op.DBDSN = cfg.DBDSN
	}
	op.EndpointURL = getEndpointURL(op)

	return op
}

func readEnv() (config, error) {
	var cfg config
	err := env.Parse(&cfg)

	return cfg, err
}

func parseFlags(op *Options) {
	flag.IntVar(&op.Storage.Interval, "i", 300, "store interval")
	flag.BoolVar(&op.Storage.Restore, "r", true, "restore from storage")
	flag.StringVar(&op.Storage.Path, "f", "/tmp/metrics-db.json", "file storage path")
	flag.StringVar(&op.DBDSN, "d", "", "db DSN")
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

func getEndpointURL(op *Options) string {
	return fmt.Sprintf("%v:%v", op.Endpoint.Hostname, op.Endpoint.Port)
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
