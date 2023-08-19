package agent

import (
	"flag"
	"fmt"
	"log"
	"net/url"
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
	return fmt.Sprintf("http://%v:%v/update/", op.ReportHost, op.ReportPort)
}

func getAddressParts(s string) (string, string, error) {
	var host, port string
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err == nil {
		host = u.Hostname()
		port = u.Port()
	}
	return host, port, err
}

func ParseFlags() *Options {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln(err)
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
			log.Fatalln(err)
		}
	}

	return op
}
