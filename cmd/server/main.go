package main

import (
	"github.com/artemmad/go-metrics-collector/internal/Server/handlers"
	"github.com/artemmad/go-metrics-collector/internal/Server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func main() {
	configFlags()
	store := storage.NewMemStorage()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handlers.MetricList(store))
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{value}", handlers.MetricCalc(store))
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", handlers.GetOneMetric(store))
	})

	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		panic(err)
	}
}
