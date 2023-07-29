package main

import (
	"net/http"

	"github.com/ilegorro/almetrics/internal/handlers"
	"github.com/ilegorro/almetrics/internal/storage"
)

func main() {
	strg := storage.NewMemStorage()
	hctx := handlers.NewUpdateHandlerContext(&strg)
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, hctx.UpdateHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
