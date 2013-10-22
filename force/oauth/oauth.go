package oauth

import (
	"fmt"
	"time"
)

const (
	loginUri     = "https://login.salesforce.com/services/oauth2/token"
	testLoginUri = "https://test.salesforce.com/services/oauth2/token"
)

type ForceOauth struct {
	AccessToken string    `json:"access_token"`
	InstanceUrl string    `json:"instance_url"`
	Id          string    `json:"id"`
	IssuedAt    time.Time `json:"issued_at"`
	Signature   string    `json:"signature"`
}

func (oauth *ForceOauth) Validate() error {
	if oauth != nil && len(oauth.InstanceUrl) != 0 && len(oauth.AccessToken) != 0 {
		return nil
	}

	return fmt.Errorf("Invalid Force Oauth Object: %#v", oauth)
}

func Authenticate(GrantType, ClientId, ClientSecret, UserName, Password, SecurityToken, Environment string) (*ForceOauth, error) {
	payload := url.Values{
		"grant_type":    {GrantType},
		"client_id":     {ClientId},
		"client_secret": {ClientSecret},
		"username":      {UserName},
		"password":      {Password + SecurityToken},
	}

	if env

}
