package server

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	ReportHost string
	ReportPort string
}

func (op *Options) GetEndpointURL() string {
	return fmt.Sprintf("%v:%v", op.ReportHost, op.ReportPort)
}

func ParseFlags() *Options {
	op := &Options{
		ReportHost: "localhost",
		ReportPort: "8080",
	}

	flag.Func("a", "host and port (default localhost:8080)", func(flagValue string) error {
		parts := strings.Split(flagValue, ":")
		if len(parts) != 2 {
			return errors.New("wrong format - host:port")
		}
		op.ReportHost = parts[0]
		op.ReportPort = parts[1]

		return nil
	})
	flag.Parse()
	return op
}
