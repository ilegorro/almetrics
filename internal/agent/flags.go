package agent

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	PollInterval   int
	ReportInterval int
	ReportHost     string
	ReportPort     string
}

func (op *Options) GetReportURL() string {
	return fmt.Sprintf("http://%v:%v/update", op.ReportHost, op.ReportPort)
}

func ParseFlags() *Options {
	op := &Options{
		ReportHost: "localhost",
		ReportPort: "8080",
	}
	flag.IntVar(&op.PollInterval, "p", 2, "poll interval")
	flag.IntVar(&op.ReportInterval, "r", 10, "report interval")
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
