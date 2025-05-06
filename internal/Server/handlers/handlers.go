package handlers

import (
	"fmt"
	"github.com/artemmad/go-metrics-collector/internal"
	"github.com/artemmad/go-metrics-collector/internal/Server/storage"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func MetricList(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder

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
