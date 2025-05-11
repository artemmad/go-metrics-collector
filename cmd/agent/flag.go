package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	serverAddress  string
	reportInterval time.Duration
	pollInterval   time.Duration
)

const (
	defaultServerAddress  = "http://localhost:8080"
	pollIntervalDefault   = 2
	reportIntervalDefault = 10

	adressEnv         = "ADDRESS"
	reportIntervalEnv = "REPORT_INTERVAL"
	pollIntervalEnv   = "POLL_INTERVAL"
)

func configFlags() {
	serverAddressEnv, serverAddressEnvExistence := os.LookupEnv(adressEnv)
	reportIntervalEnv, reportIntervalEnvExistence := os.LookupEnv(reportIntervalEnv)
	pollIntervalEnv, pollIntervalEnvExistence := os.LookupEnv(pollIntervalEnv)

	serverAddressParam := flag.String("a", defaultServerAddress, "The address to bind the server to, ex. http://localhost:8080")
	reportIntervalIntParam := flag.Int("r", reportIntervalDefault, "The interval in seconds between send of metrics to the server")
	pollIntervalIntParam := flag.Int("p", pollIntervalDefault, "The interval between scrap of metrics in seconds")

	flag.Parse()

	var resServerAddress string
	var resReportInterval int
	var resPollInterval int

	if serverAddressEnvExistence {
		resServerAddress = serverAddressEnv
	} else {
		resServerAddress = *serverAddressParam
	}

	if reportIntervalEnvExistence {
		v, err := strconv.ParseInt(reportIntervalEnv, 10, 0)
		if err != nil {
			panic(err)
		}
		resReportInterval = int(v)
	} else {
		resReportInterval = *reportIntervalIntParam
	}

	if pollIntervalEnvExistence {
		v, err := strconv.ParseInt(pollIntervalEnv, 10, 0)
		if err != nil {
			panic(err)
		}
		resPollInterval = int(v)
	} else {
		resPollInterval = *pollIntervalIntParam
	}

	if !strings.HasPrefix(resServerAddress, "http://") && !strings.HasPrefix(resServerAddress, "https://") {
		serverAddress = "http://" + resServerAddress
	} else {
		serverAddress = resServerAddress
	}
	reportInterval = time.Duration(resReportInterval) * time.Second
	pollInterval = time.Duration(resPollInterval) * time.Second
}
