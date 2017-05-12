package force

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	grantType    = "password"
	loginUri     = "https://login.salesforce.com/services/oauth2/token"
	testLoginUri = "https://test.salesforce.com/services/oauth2/token"

	invalidSessionErrorCode = "INVALID_SESSION_ID"
)

type forceOauth struct {
	AccessToken string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
	Id          string `json:"id"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`

	clientId      string
	clientSecret  string
	refreshToken  string
	userName      string
	password      string
	securityToken string
	environment   string
}

func (oauth *forceOauth) Validate() error {
	if oauth == nil || len(oauth.InstanceUrl) == 0 || len(oauth.AccessToken) == 0 {
		return fmt.Errorf("Invalid Force Oauth Object: %#v", oauth)
	}

	return nil
}

func (oauth *forceOauth) Expired(apiErrors ApiErrors) bool {
	for _, err := range apiErrors {
		if err.ErrorCode == invalidSessionErrorCode {
			return true
		}
	}

	return false
}

func (oauth *forceOauth) Authenticate() error {
	payload := url.Values{
		"grant_type":    {grantType},
		"client_id":     {oauth.clientId},
		"client_secret": {oauth.clientSecret},
		"username":      {oauth.userName},
		"password":      {fmt.Sprintf("%v%v", oauth.password, oauth.securityToken)},
	}

	// Build Uri
	uri := loginUri
	if oauth.environment == "sandbox" {
		uri = testLoginUri
	}

	// Build Body
	body := strings.NewReader(payload.Encode())

	// Build Request
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return fmt.Errorf("Error creating authentication request: %v", err)
	}

	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", responseType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending authentication request: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading authentication response bytes: %v", err)
	}

	// Attempt to parse response as a force.com api error
	apiError := &ApiError{}
	if err := json.Unmarshal(respBytes, apiError); err == nil {
		// Check if api error is valid
		if apiError.Validate() {
			return apiError
		}
	}

	if err := json.Unmarshal(respBytes, oauth); err != nil {
		return fmt.Errorf("Unable to unmarshal authentication response: %v", err)
	}

	return nil
}
