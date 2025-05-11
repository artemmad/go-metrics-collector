package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func Test_configFlags_PriorityEnvOverFlag(t *testing.T) {
	t.Setenv(Adress_env, "env-host:8888")

	resetFlags()
	os.Args = []string{"cmd", "-a=cli-host:7777"}

	configFlags()

	assert.Equal(t, "env-host:8888", serverAddress)
}

func Test_configFlags_CLIUsedIfNoEnv(t *testing.T) {
	os.Unsetenv(Adress_env)

	resetFlags()
	os.Args = []string{"cmd", "-a=cli-host:5555"}

	configFlags()

	assert.Equal(t, "cli-host:5555", serverAddress)
}

func Test_configFlags_DefaultUsedIfNoEnvAndFlag(t *testing.T) {
	os.Unsetenv(Adress_env)

	resetFlags()
	os.Args = []string{"cmd"}

	configFlags()

	assert.Equal(t, defaultServerAddress, serverAddress)
}
