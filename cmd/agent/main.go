package main

import (
	agent "github.com/artemmad/go-metrics-collector/internal/Agent"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	go func() {
		for {
			agent.UpdateMetrics()
			time.Sleep(pollInterval)
		}
	}()
	for {
		agent.ReportMetrics()
		time.Sleep(reportInterval)
	}
}
