package security

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"
	"io/ioutil"
	"net/url"
	"strings"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type EndPoints struct {
	AuthEndpoint  string `json:"authorization_endpoint"`
	TokenEndpoint string `json:"token_endpoint"`
}

var cfTokens *Tokens

//
// get the Authorization and Token endpoints from API endpoint
//
func GetEndPoints(api string) (*EndPoints, error) {
	url := api + CF_INFO_PATH
	resp, err := DoRequest("GET", url, "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("Error get auth and token endpoints from " + url + ": " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error get endpoints: failed to read response body ")
	}

	var endpoints EndPoints
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return nil, errors.New("Error get endpoints: failed to unmarshall json body")
	}

	return &endpoints, nil
}

//
// Get Access/Refresh Tokens from login server
//

func Login(c *config.CfConfig) error {
	cfTokens = &Tokens{}

	endpoints, err := GetEndPoints(c.Api)
	if err != nil {
		return err
	}

	authUrl := endpoints.AuthEndpoint + CF_AUTH_PATH
	grantType := strings.ToLower(c.GrantType)

	var form url.Values
	if grantType == "password" {
		form = url.Values{
			"grant_type": {"password"},
			"username":   {c.User},
			"password":   {c.Pass},
		}
	} else if grantType == "client_credentials" {
		form = url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {c.ClientId},
			"client_secret": {c.Secret},
		}
	} else {
		return errors.New("Not supported grant type :" + grantType)
	}

	headers := map[string]string{}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["charset"] = "utf-8"

	var token string
	if grantType == "password" {
		token = "Basic Y2Y6"
	} else {
		token = c.ClientId + ":" + c.Secret
		token = "Basic " + base64.StdEncoding.EncodeToString([]byte(token))
	}

	resp, err := DoRequest("POST", authUrl, token, headers, strings.NewReader(form.Encode()))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Error get oauth tokens: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error get oauth tokens: failed to read auth response body " + err.Error())
	}

	err = json.Unmarshal(body, cfTokens)
	return err
}

func GetOAuthToken() string {
	return cfTokens.AccessToken
}
