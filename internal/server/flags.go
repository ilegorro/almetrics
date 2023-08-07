package server

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

type Options struct {
	ReportHost string
	ReportPort string
}

func (op *Options) GetEndpointURL() string {
	return fmt.Sprintf("%v:%v", op.ReportHost, op.ReportPort)
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
	flag.Func("a", "host and port (default localhost:8080)", func(flagValue string) error {
		op.ReportHost, op.ReportPort, err = getAddressParts(flagValue)
		if err != nil {
			return err
		}
		return nil
	})
	flag.Parse()

	if cfg.Address != "" {
		op.ReportHost, op.ReportPort, err = getAddressParts(cfg.Address)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return op
}
