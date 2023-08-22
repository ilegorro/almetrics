package handlers

import (
	"log"

	"github.com/ilegorro/almetrics/internal/common"
)

type HandlerContext struct {
	strg     common.Repository
	syncPath string
}

func NewHandlerContext(strg common.Repository, syncPath string) *HandlerContext {
	if strg == nil {
		log.Fatalln("Storage is not defined")
	}
	return &HandlerContext{strg: strg, syncPath: syncPath}
}
