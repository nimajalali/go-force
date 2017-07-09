package force

import (
	"testing"

	"github.com/goguardian/go-force/sobjects"
)

func TestCreateWithAccessToken(t *testing.T) {

	// Manually grab an OAuth token, so that we can pass it into CreateWithAccessToken
	oauth := &forceOauth{
		clientID:      testClientID,
		clientSecret:  testClientSecret,
		userName:      testUserName,
		password:      testPassword,
		securityToken: testSecurityToken,
		environment:   testEnvironment,
	}

	forceAPI := &ForceAPI{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	err := forceAPI.oauth.Authenticate()
	if err != nil {
		t.Fatalf("Unable to authenticate: %#v", err)
	}
	if err := forceAPI.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// We shouldn't hit any errors creating a new force instance and manually passing in these oauth details now.
	newForceAPI, err := CreateWithAccessToken(testVersion, testClientID, forceAPI.oauth.AccessToken, forceAPI.oauth.InstanceUrl)
	if err != nil {
		t.Fatalf("Unable to create new force api instance using pre-defined oauth details: %#v", err)
	}
	if err := newForceAPI.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// We should be able to make a basic query now with the newly created object (i.e. the oauth details should be correctly usable).
	_, err = newForceAPI.DescribeSObject(&sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}
}
