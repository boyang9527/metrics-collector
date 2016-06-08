package main

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/security"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/server"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"
	"os"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "configuration file")
	flag.Parse()

	var conf *config.Config
	var err error

	if path == "" {
		conf = config.DefaultConfig()
	} else {
		conf, err = config.LoadConfigFromFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read config file: %s\n", err.Error())
			os.Exit(1)
		}
	}

	fmt.Println(conf.ToString())

	err = InitailizeLogger(&conf.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to intialize logger: %s\n", err.Error())
		os.Exit(1)
	}

	err = security.Login(&conf.Cf)
	if err != nil {
		Logger.Error("failed-to-login-cloudfoundry", err)
		os.Exit(1)
	}

	s := server.NewServer(&conf.Server)
	s.Start()
}
