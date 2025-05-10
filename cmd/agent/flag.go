package main

import (
	"flag"
	"time"
)

var (
	serverAddress  string
	reportInterval time.Duration
	pollInterval   time.Duration
)

const (
	pollIntervalDefault   = 2 * time.Second
	reportIntervalDefault = 10 * time.Second
)

func configFlags() {
	flag.StringVar(&serverAddress, "a", ":8080", "The address to bind the server to, ex. localhost:8080")
	flag.DurationVar(&reportInterval, "r", reportIntervalDefault, "The interval between send of metrics to the server")
	flag.DurationVar(&pollInterval, "p", pollIntervalDefault, "The interval between scrap of metrics")

	flag.Parse()
}
