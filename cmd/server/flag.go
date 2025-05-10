package main

import (
	"flag"
)

var (
	serverAddress string
)

func configFlags() {
	flag.StringVar(&serverAddress, "a", ":8080", "The address to bind the server to, ex. localhost:8080")

	flag.Parse()
}
