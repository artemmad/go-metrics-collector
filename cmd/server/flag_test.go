package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_configFlags(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-a=127.0.0.1:9000"}

	configFlags()

	assert.Equal(t, "127.0.0.1:9000", serverAddress)
}
