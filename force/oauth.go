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
	loginUri     = "https://login.salesforce.com/services/oauth2/token"
	testLoginUri = "https://test.salesforce.com/services/oauth2/token"
)

type ForceOauth struct {
	AccessToken string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
	Id          string `json:"id"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

func (oauth *ForceOauth) Validate() error {
	if oauth == nil || len(oauth.InstanceUrl) == 0 || len(oauth.AccessToken) == 0 {
		return fmt.Errorf("Invalid Force Oauth Object: %#v", oauth)
	}

	return nil
}

func authenticate(GrantType, ClientId, ClientSecret, UserName, Password, SecurityToken, Environment string) (*ForceOauth, error) {
	payload := url.Values{
		"grant_type":    {GrantType},
		"client_id":     {ClientId},
		"client_secret": {ClientSecret},
		"username":      {UserName},
		"password":      {Password + SecurityToken},
	}

	// Build Uri
	uri := loginUri
	if Environment == "sandbox" {
		uri = testLoginUri
	}

	// Build Body
	body := strings.NewReader(payload.Encode())

	// Build Request
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, fmt.Errorf("Error creating authenitcation request: %v", err)
	}

	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", responseType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending authentication request: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading authentication response bytes: %v", err)
	}

	// Attempt to parse response as a force.com api error
	apiError := &ApiError{}
	if err := json.Unmarshal(respBytes, apiError); err == nil {
		// Check if api error is valid
		if len(apiError.Fields) != 0 || len(apiError.Message) != 0 || len(apiError.ErrorCode) != 0 || len(apiError.ErrorName) != 0 || len(apiError.ErrorDescription) != 0 {
			return nil, apiError
		}
	}

	// Attempt to parse response into ForceOauth object
	respObject := &ForceOauth{}
	if err := json.Unmarshal(respBytes, respObject); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal authentication response: %v", err)
	}

	fmt.Println(respObject)

	return respObject, nil
}
