package util

import (
	"fmt"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	"github.com/pivotal-golang/lager"
	"os"
	"path/filepath"
)

var Logger = lager.NewLogger("as-metrics-collector")

func getLogLevel(level string) lager.LogLevel {
	switch level {
	case "DEBUG":
		return lager.DEBUG
	case "INFO":
		return lager.INFO
	case "ERROR":
		return lager.ERROR
	case "FATAL":
		return lager.FATAL

	default:
		return lager.INFO
	}
}

func InitailizeLogger(c *config.LoggingConfig) {

	logLevel := getLogLevel(c.Level)

	if c.LogToStdout {
		Logger.RegisterSink(lager.NewWriterSink(os.Stdout, logLevel))
	}

	if c.File != "" {
		var file *os.File

		info, err := os.Stat(c.File)
		if err == nil {
			if info.IsDir() {
				fmt.Fprintf(os.Stderr, "log file '%s' is a directory\n", c.File)
				os.Exit(1)
			}

			file, err = os.OpenFile(c.File, os.O_APPEND|os.O_RDWR, 0644)

		} else {

			err = os.MkdirAll(filepath.Dir(c.File), 0744)
			if err == nil {
				file, err = os.Create(c.File)
			}
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log file %s \n", err.Error())
			os.Exit(1)
		}
		Logger.RegisterSink(lager.NewWriterSink(file, logLevel))
	}
}
