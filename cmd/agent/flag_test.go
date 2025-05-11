package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func Test_configFlags_Priority(t *testing.T) {
	t.Setenv(Adress_env, "https://env-address")
	t.Setenv(ReportInterval_env, "99")
	t.Setenv(PollInterval_env, "77")

	resetFlags()
	os.Args = []string{
		"cmd",
		"-a=cli-address",
		"-r=33",
		"-p=22",
	}

	configFlags()

	assert.Equal(t, "https://env-address", serverAddress)
	assert.Equal(t, 99*time.Second, reportInterval)
	assert.Equal(t, 77*time.Second, pollInterval)
}

func Test_configFlags_CLIOnly(t *testing.T) {
	os.Unsetenv(Adress_env)
	os.Unsetenv(ReportInterval_env)
	os.Unsetenv(PollInterval_env)

	resetFlags()
	os.Args = []string{
		"cmd",
		"-a=myhost:5000",
		"-r=15",
		"-p=5",
	}

	configFlags()

	assert.Equal(t, "http://myhost:5000", serverAddress)
	assert.Equal(t, 15*time.Second, reportInterval)
	assert.Equal(t, 5*time.Second, pollInterval)
}

func Test_configFlags_Defaults(t *testing.T) {
	os.Unsetenv(Adress_env)
	os.Unsetenv(ReportInterval_env)
	os.Unsetenv(PollInterval_env)

	resetFlags()
	os.Args = []string{"cmd"}

	configFlags()

	assert.Equal(t, defaultServerAddress, serverAddress)
	assert.Equal(t, time.Duration(reportIntervalDefault)*time.Second, reportInterval)
	assert.Equal(t, time.Duration(pollIntervalDefault)*time.Second, pollInterval)
}
