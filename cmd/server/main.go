package main

import (
	"fmt"
	"github.com/artemmad/go-metrics-collector/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

const (
	gaugeType   = "gauge"
	counterType = "counter"
)

func main() {
	store := storage.NewMemStorage()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/", metricCalc(store))
	mux.HandleFunc("GET /", metricList(store))
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "404 page not found", http.StatusNotFound)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func metricList(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder

		b.WriteString("GAUGE:\n")
		for k, v := range store.GetGauges() {
			b.WriteString(fmt.Sprintf("\t%s: %f\n", k, v))
		}

		b.WriteString("COUNTER:\n")
		for k, v := range store.GetCounters() {
			b.WriteString(fmt.Sprintf("\t%s: %d\n", k, v))
		}

		w.Write([]byte(b.String()))
	}
}

func metricCalc(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 || parts[0] != "update" {
			http.Error(w, "invalid URL format", http.StatusBadRequest)
			return
		}

		metricType := strings.ToLower(parts[1])
		name := strings.TrimSpace(parts[2])
		valueStr := strings.TrimSpace(parts[3])

		switch metricType {
		case gaugeType:
			val, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				http.Error(w, "invalid gauge value", http.StatusBadRequest)
				return
			}
			store.SetGauge(name, val)

		case counterType:
			val, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid counter value", http.StatusBadRequest)
				return
			}
			store.SetCounter(name, val)

		default:
			http.Error(w, "unknown metric type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
