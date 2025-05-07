package handlers

import (
	"fmt"
	"github.com/artemmad/go-metrics-collector/internal"
	"github.com/artemmad/go-metrics-collector/internal/Server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func MetricList(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		b.WriteString("GAUGE:\n")
		gauges := store.GetGauges()
		gaugeKeys := make([]string, 0, len(gauges))
		for k := range gauges {
			gaugeKeys = append(gaugeKeys, k)
		}
		sort.Strings(gaugeKeys)
		for _, k := range gaugeKeys {
			b.WriteString(fmt.Sprintf("\t%s: %f\n", k, gauges[k]))
		}

		b.WriteString("COUNTER:\n")
		counters := store.GetCounters()
		counterKeys := make([]string, 0, len(counters))
		for k := range counters {
			counterKeys = append(counterKeys, k)
		}
		sort.Strings(counterKeys)
		for _, k := range counterKeys {
			b.WriteString(fmt.Sprintf("\t%s: %d\n", k, counters[k]))
		}

		w.Write([]byte(b.String()))
		w.WriteHeader(http.StatusOK)
	}
}

func MetricCalc(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := strings.ToLower(chi.URLParam(r, "metricType"))
		name := strings.TrimSpace(chi.URLParam(r, "metricName"))
		valueStr := strings.TrimSpace(chi.URLParam(r, "value"))

		switch metricType {
		case internal.GaugeType:
			val, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				http.Error(w, "invalid gauge value", http.StatusBadRequest)
				return
			}
			store.SetGauge(name, val)

		case internal.CounterType:
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

func GetOneMetric(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		w.Header().Add("Content-Type", "text/plain")
		switch metricType {
		case internal.GaugeType:
			res, ok := store.GetGauges()[metricName]
			if !ok {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			w.Write([]byte(strconv.FormatFloat(res, 'f', -1, 64)))
		case internal.CounterType:
			res, ok := store.GetCounters()[metricName]
			if !ok {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			w.Write([]byte(strconv.FormatInt(res, 10)))
		default:
			http.Error(w, "unknown metric type", http.StatusNotFound)
		}
		w.WriteHeader(http.StatusOK)
	}
}
