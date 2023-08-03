package agent

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type Options struct {
	PollInterval   int
	ReportInterval int
	ReportHost     string
	ReportPort     string
}

func (op *Options) GetReportURL() string {
	return fmt.Sprintf("http://%v:%v/update", op.ReportHost, op.ReportPort)
}

func getAddressParts(s string) (host string, port string, err error) {
	parts := strings.Split(s, ":")
	if len(parts) == 2 {
		host = parts[0]
		port = parts[1]
	} else {
		err = errors.New("wrong format - host:port")
	}
	return
}

func ParseFlags() *Options {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	op := &Options{
		ReportHost: "localhost",
		ReportPort: "8080",
	}
	flag.IntVar(&op.PollInterval, "p", 2, "poll interval")
	flag.IntVar(&op.ReportInterval, "r", 10, "report interval")
	flag.Func("a", "host and port (default localhost:8080)", func(flagValue string) error {
		op.ReportHost, op.ReportPort, err = getAddressParts(flagValue)
		if err != nil {
			return err
		}
		return nil
	})
	flag.Parse()

	if cfg.PollInterval != 0 {
		op.PollInterval = cfg.PollInterval
	}
	if cfg.ReportInterval != 0 {
		op.ReportInterval = cfg.ReportInterval
	}
	if cfg.Address != "" {
		op.ReportHost, op.ReportPort, err = getAddressParts(cfg.Address)
		if err != nil {
			panic(err)
		}
	}

	return op
}
