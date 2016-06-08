package security_test

import (
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/fakes"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestSecurity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Security Suite")
}

var cfServer = &fakes.FakeCfServer{}

var _ = BeforeSuite(func() {

	testLoggingConfig := config.LoggingConfig{
		Level:       "info",
		File:        "",
		LogToStdout: true,
	}

	InitailizeLogger(&testLoggingConfig)
	go cfServer.Start()

	// allow the server to be ready
	time.Sleep(2 * time.Second)
})

var _ = AfterSuite(func() {
	cfServer.Stop()
})
