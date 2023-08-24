package handlers

import (
	"github.com/ilegorro/almetrics/internal/common"
)

type HandlerContext struct {
	strg     common.Repository
	syncPath string
}

func NewHandlerContext(strg common.Repository, syncPath string) *HandlerContext {
	logger := common.SugaredLogger()
	if strg == nil {
		logger.Fatalln("Storage is not defined")
	}
	return &HandlerContext{strg: strg, syncPath: syncPath}
}
