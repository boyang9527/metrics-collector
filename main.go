package main

import (
	"flag"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/security"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/server"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"
	"os"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.Parse()

	var conf *config.Config
	if configFile == "" {
		conf = config.DefaultConfig()
	} else {
		conf = config.LoadConfigFromFile(configFile)
	}

	InitailizeLogger(&conf.Logging)

	err := security.Login(&conf.Cf)
	if err != nil {
		Logger.Error("failed-to-login-cloudfoundry", err)
		os.Exit(1)
	}

	s := server.NewServer(&conf.Server)
	s.Start()
}
