package main

import (
	"flag"
	"strings"
	"time"
)

var (
	serverAddress  string
	reportInterval time.Duration
	pollInterval   time.Duration
)

const (
	pollIntervalDefault   = 2
	reportIntervalDefault = 10
)

func configFlags() {
	flag.StringVar(&serverAddress, "a", "http://localhost:8080", "The address to bind the server to, ex. http://localhost:8080")
	reportIntervalInt := flag.Int("r", reportIntervalDefault, "The interval in seconds between send of metrics to the server")
	pollIntervalInt := flag.Int("p", pollIntervalDefault, "The interval between scrap of metrics in seconds")

	flag.Parse()

	if !strings.HasPrefix(serverAddress, "http://") && !strings.HasPrefix(serverAddress, "https://") {
		serverAddress = "http://" + serverAddress
	}
	reportInterval = time.Duration(*reportIntervalInt) * time.Second
	pollInterval = time.Duration(*pollIntervalInt) * time.Second
}
