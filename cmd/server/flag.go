package main

import (
	"flag"
	"os"
)

var (
	serverAddress string
)

const (
	defaultServerAddress = ":8080"

	adressEnv = "ADDRESS"
)

func configFlags() {
	serverAddressEnv, serverAdressEnvExistence := os.LookupEnv(adressEnv)
	serverAdressParam := flag.String("a", defaultServerAddress, "The address to bind the server to, ex. localhost:8080")
	flag.Parse()

	if serverAdressEnvExistence {
		serverAddress = serverAddressEnv
	} else {
		serverAddress = *serverAdressParam
	}
}
