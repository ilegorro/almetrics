package handlers

import "github.com/ilegorro/almetrics/internal/storage"

type HandlerContext struct {
	strg storage.Repository
}

func NewHandlerContext(strg storage.Repository) *HandlerContext {
	if strg == nil {
		panic("Storage is not defined")
	}
	return &HandlerContext{strg: strg}
}
