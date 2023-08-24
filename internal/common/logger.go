package common

import "go.uber.org/zap"

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		SugaredLogger().Fatalf("Unable to set up logger: %+v", err)
	}
	zap.ReplaceGlobals(l)
}

func Logger() *zap.Logger {
	return zap.L()
}

func SugaredLogger() *zap.SugaredLogger {
	return zap.L().Sugar()
}
