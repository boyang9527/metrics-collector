package fakes

import (
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"

	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
)

const (
	FAKE_CF_SERVER_ADDR = "127.0.0.1:8989"
	FAKE_AUTH_ENDPOINT  = "http://127.0.0.1:8989"
	FAKE_TOKEN_ENDPOINT = "http://127.0.0.1:8989"
	FAKE_OAUTH_TOKEN    = "fake-oauth-token"
	FAKE_REFRESH_TOKEN  = "fake-refresh-token"
)

type FakeCfServer struct {
	listener *net.Listener
}

var FakeCfConfig = config.CfConfig{
	Api:       "http://" + FAKE_CF_SERVER_ADDR,
	GrantType: "password",
	User:      "fake-user",
	Pass:      "fake-pass",
	ClientId:  "fake-client",
	Secret:    "fake-secret",
}

var infoBody = []byte(`
{
   "name": "",
   "build": "",
   "support": "http://support.cloudfoundry.com",
   "version": 0,
   "description": "",
   "authorization_endpoint": "{AUTH_ENDPOINT}",
   "token_endpoint": "{TOKEN_ENDPOINT}",
   "min_cli_version": null,
   "min_recommended_cli_version": null,
   "api_version": "2.48.0",
   "app_ssh_endpoint": "ssh.bosh-lite.com:2222",
   "app_ssh_host_key_fingerprint": "a6:d1:08:0b:b0:cb:9b:5f:c4:ba:44:2a:97:26:19:8a",
   "app_ssh_oauth_client": "ssh-proxy",
   "routing_endpoint": "https://api.bosh-lite.com/routing",
   "logging_endpoint": "wss://loggregator.bosh-lite.com:443",
   "doppler_logging_endpoint": "wss://doppler.bosh-lite.com:4443",
   "user": "38b2f682-04bf-48af-9e08-0325aa5c4ea9"
}
`)

var authBody = []byte(`
{
	"access_token":"{OAUTH_TOKEN}",
	"token_type":"bearer",
	"refresh_token":"{REFRESH_TOKEN}",
	"expires_in":43199,
	"scope":"openid cloud_controller.read password.write cloud_controller.write",
	"jti":"a735f90f-0b49-447d-8f9d-ae2fbc1491dd"}				
`)

func (fake *FakeCfServer) Start() {
	listener, err := net.Listen("tcp", FAKE_CF_SERVER_ADDR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start cf mock server: %s\n", err.Error())
		os.Exit(1)
	}
	fake.listener = &listener

	r := mux.NewRouter()
	r.Methods("GET").Path(CF_INFO_PATH).HandlerFunc(handleInfo)
	r.Methods("POST").Path(CF_AUTH_PATH).HandlerFunc(handleLogin)

	http.Handle("/", r)
	http.Serve(listener, nil)
}

func (fake *FakeCfServer) Stop() {
	(*fake.listener).Close()
}

func handleInfo(w http.ResponseWriter, r *http.Request) {

	b := bytes.Replace(infoBody, []byte("{AUTH_ENDPOINT}"), []byte(FAKE_AUTH_ENDPOINT), -1)
	b = bytes.Replace(b, []byte("{TOKEN_ENDPOINT}"), []byte(FAKE_TOKEN_ENDPOINT), -1)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")
	if grantType != "password" && grantType != "client_credentials" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(CreateJsonErrorResponse("Error-login-CF", "invalid grant_type"))
		return
	}

	authHeader := r.Header.Get("Authorization")

	if grantType == "password" {
		if authHeader != "Basic Y2Y6" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write(CreateJsonErrorResponse("Error-Get-login-CF", "invalid authorization header"))
			return
		}

		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user != FakeCfConfig.User || pass != FakeCfConfig.Pass {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write(CreateJsonErrorResponse("Error-Get-login-CF", "invalid login credentials"))
			return
		}
	} else {
		token := "Basic " + base64.StdEncoding.EncodeToString([]byte(FakeCfConfig.ClientId+":"+FakeCfConfig.Secret))
		if authHeader != token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write(CreateJsonErrorResponse("Error-Get-login-CF", "invalid authorization header"))
			return
		}

		clientId := r.FormValue("client_id")
		secret := r.FormValue("client_secret")

		if clientId != FakeCfConfig.ClientId || secret != FakeCfConfig.Secret {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write(CreateJsonErrorResponse("Error-Get-login-CF", "invalid client credentials"))
			return
		}

	}

	b := bytes.Replace(authBody, []byte("{OAUTH_TOKEN}"), []byte(FAKE_OAUTH_TOKEN), -1)
	b = bytes.Replace(b, []byte("{REFRESH_TOKEN}"), []byte(FAKE_REFRESH_TOKEN), -1)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
