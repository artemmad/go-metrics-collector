package internal

import (
	"fmt"
	internal2 "github.com/artemmad/go-metrics-collector/internal"
	"math/rand/v2"
	"net/http"
	"runtime"
)

var (
	randomValue   float64
	pollCount     int64
	ServerAddress = "http://localhost:8080"
)

func SetServerAddress(serverAddress string) {
	ServerAddress = serverAddress
}

func ReportMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ReportGaugeMetric("Alloc", float64(m.Alloc))
	ReportGaugeMetric("BuckHashSys", float64(m.BuckHashSys))
	ReportGaugeMetric("Frees", float64(m.Frees))
	ReportGaugeMetric("GCCPUFraction", m.GCCPUFraction)
	ReportGaugeMetric("GCSys", float64(m.GCSys))
	ReportGaugeMetric("HeapAlloc", float64(m.HeapAlloc))
	ReportGaugeMetric("HeapIdle", float64(m.HeapIdle))
	ReportGaugeMetric("HeapInuse", float64(m.HeapInuse))
	ReportGaugeMetric("HeapReleased", float64(m.HeapReleased))
	ReportGaugeMetric("HeapSys", float64(m.HeapSys))
	ReportGaugeMetric("LastGC", float64(m.LastGC))
	ReportGaugeMetric("Lookups", float64(m.Lookups))
	ReportGaugeMetric("MCacheInuse", float64(m.MCacheInuse))
	ReportGaugeMetric("MCacheSys", float64(m.MCacheSys))
	ReportGaugeMetric("MSpanInuse", float64(m.MSpanInuse))
	ReportGaugeMetric("MSpanSys", float64(m.MSpanSys))
	ReportGaugeMetric("Mallocs", float64(m.Mallocs))
	ReportGaugeMetric("NextGC", float64(m.NextGC))
	ReportGaugeMetric("NumForcedGC", float64(m.NumForcedGC))
	ReportGaugeMetric("NumGC", float64(m.NumGC))
	ReportGaugeMetric("OtherSys", float64(m.OtherSys))
	ReportGaugeMetric("PauseTotalNs", float64(m.PauseTotalNs))
	ReportGaugeMetric("StackInuse", float64(m.StackInuse))
	ReportGaugeMetric("StackSys", float64(m.StackSys))
	ReportGaugeMetric("Sys", float64(m.Sys))
	ReportGaugeMetric("TotalAlloc", float64(m.TotalAlloc))

	ReportCounterMetric("PollCount", pollCount)
	ReportGaugeMetric("RandomValue", randomValue)
}

func ReportGaugeMetric(name string, value float64) {
	ReportMetric(internal2.GaugeType, name, value)
}

func ReportCounterMetric(name string, value int64) {
	ReportMetric(internal2.CounterType, name, value)
}

func ReportMetric(metricType string, name string, val interface{}) {
	url := fmt.Sprintf("%s/update/%s/%s/%v", ServerAddress, metricType, name, val)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

func UpdateMetrics() {
	pollCount++
	randomValue = rand.Float64()
}
