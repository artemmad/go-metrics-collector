package handlers

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

func MetricList(store storage.Storage) http.HandlerFunc {
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

func MetricCalc(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 || parts[0] != "update" {
			http.Error(w, "invalid URL format", http.StatusNotFound)
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
