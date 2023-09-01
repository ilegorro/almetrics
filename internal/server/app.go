package server

import (
	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/server/config"
)

type App struct {
	strg            common.Repository
	options         *config.Options
	syncFileStorage bool
}

func NewApp(strg common.Repository, op *config.Options) *App {
	logger := common.SugaredLogger()
	if strg == nil {
		logger.Fatalln("Storage is not defined")
	}
	syncFS := (op.Storage.Interval == 0 && op.Storage.Path != "")

	return &App{strg: strg, options: op, syncFileStorage: syncFS}
}
