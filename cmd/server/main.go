package main

import (
	"github.com/artemmad/go-metrics-collector/internal/Server/handlers"
	"github.com/artemmad/go-metrics-collector/internal/Server/storage"
	"net/http"
)

func main() {
	store := storage.NewMemStorage()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/", handlers.MetricCalc(store))
	mux.HandleFunc("GET /", handlers.MetricList(store))
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "404 page not found", http.StatusNotFound)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
