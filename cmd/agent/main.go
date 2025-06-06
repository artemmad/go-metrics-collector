package main

import (
	agent "github.com/artemmad/go-metrics-collector/internal/Agent"
	"time"
)

func main() {
	configFlags()
	agent.SetServerAddress(serverAddress)
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
