package security_test

import (
	"fmt"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/fakes"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/security"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	Describe("get CF endpoints", func() {
		var (
			api       string
			endpoints *EndPoints
			err       error
		)

		JustBeforeEach(func() {
			endpoints, err = GetEndPoints(api)
		})

		Context("when API endpoint is reachable", func() {
			BeforeEach(func() {
				api = fakes.FakeCfConfig.Api
			})

			It("should return correct auth and token endpoints", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(endpoints).NotTo(BeNil())
				Expect(endpoints.AuthEndpoint).To(Equal(fakes.FAKE_AUTH_ENDPOINT))
				Expect(endpoints.TokenEndpoint).To(Equal(fakes.FAKE_TOKEN_ENDPOINT))
			})
		})

		Context("when API endpoint is not reachable", func() {
			BeforeEach(func() {
				api = "http://www.not-exist-server.com"
			})

			It("should error and return nil endpoint", func() {
				Expect(err).To(HaveOccurred())
				Expect(endpoints).To(BeNil())
			})
		})

		Context("when API endpoint is not serving the given path ", func() {
			BeforeEach(func() {
				api = fakes.FakeCfConfig.Api + "/not-exist-path"
			})

			It("should error and return nil endpoint", func() {
				Expect(err).To(HaveOccurred())
				Expect(endpoints).To(BeNil())
			})
		})
	})

	Describe("Login CF to get tokens", func() {
		var (
			conf config.CfConfig
			err  error
		)

		JustBeforeEach(func() {
			err = Login(&conf)
		})

		Context("when login with password", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
			})

			It("should get correct oauth token", func() {
				Expect(err).To(BeNil())
				Expect(GetOAuthToken()).To(Equal(fakes.FAKE_OAUTH_TOKEN))
			})
		})

		Context("when login with client credential", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.GrantType = "client_credentials"
			})

			It("should get correct oauth token", func() {
				Expect(err).To(BeNil())
				Expect(GetOAuthToken()).To(Equal(fakes.FAKE_OAUTH_TOKEN))
			})
		})

		Context("when login with wrong api endpoint", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.Api = "not-exist-api"
			})

			It("should error and return empty oauth token", func() {
				Expect(err).NotTo(BeNil())
				Expect(GetOAuthToken()).To(BeEmpty())
			})
		})

		Context("when login with not supported grant type", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.GrantType = "not-exist-type"
			})

			It("should error and return empty oauth token", func() {
				Expect(err).NotTo(BeNil())
				Expect(GetOAuthToken()).To(BeEmpty())
			})
		})

		Context("when login with wrong password", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.Pass = "not-exist-password"
			})

			It("should error and return empty oauth token", func() {
				Expect(err).NotTo(BeNil())
				Expect(GetOAuthToken()).To(BeEmpty())
			})
		})

		Context("when login with wrong user", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.User = "not-exist-user"
			})

			It("should error and return empty oauth token", func() {
				Expect(err).NotTo(BeNil())
				Expect(GetOAuthToken()).To(BeEmpty())
			})
		})

		Context("when login with wrong client id", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.GrantType = "client_credentials"
				conf.ClientId = "not-exist-client-id"
			})

			It("should error and return empty oauth token", func() {
				Expect(err).NotTo(BeNil())
				Expect(GetOAuthToken()).To(BeEmpty())
			})
		})

		Context("when login with wrong client secret", func() {
			BeforeEach(func() {
				conf = fakes.FakeCfConfig
				conf.GrantType = "client_credentials"
				conf.ClientId = "not-exist-client-secret"
			})

			It("should error and return empty oauth token", func() {
				Expect(err).NotTo(BeNil())
				Expect(GetOAuthToken()).To(BeEmpty())
			})
		})

	})

})
