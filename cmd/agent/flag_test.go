package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_configFlags_withDurations(t *testing.T) {
	// Сброс флагов перед тестом
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Подставим аргументы
	os.Args = []string{
		"cmd",
		"-a=127.0.0.1:9000",
		"-r=15",
		"-p=3",
	}

	// Вызов конфигурации
	configFlags()

	// Проверки
	assert.Equal(t, "http://127.0.0.1:9000", serverAddress)
	assert.Equal(t, 15*time.Second, reportInterval)
	assert.Equal(t, 3*time.Second, pollInterval)
}
