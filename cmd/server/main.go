package main

import (
	"log"
	"net/http"

	"github.com/ilegorro/almetrics/internal/handlers"
	"github.com/ilegorro/almetrics/internal/server"
	"github.com/ilegorro/almetrics/internal/storage"
)

func main() {
	op := server.ParseFlags()
	strg := storage.NewMemStorage()
	hctx := handlers.NewHandlerContext(strg)
	router := server.MetricsRouter(hctx)
	endPoint := op.GetEndpointURL()

	if err := http.ListenAndServe(endPoint, router); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
