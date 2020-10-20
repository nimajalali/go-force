package force

import (
	"context"
	"testing"

	"github.com/EverlongProject/go-force/sobjects"
)

func TestCreateWithAccessToken(t *testing.T) {

	// Manually grab an OAuth token, so that we can pass it into CreateWithAccessToken
	oauth := &forceOauth{
		clientId:      testClientId,
		clientSecret:  testClientSecret,
		userName:      testUserName,
		password:      testPassword,
		securityToken: testSecurityToken,
		environment:   testEnvironment,
		loginUri:      testLoginUri,
	}

	forceApi := &ForceApi{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}
	ctx := context.Background()
	err := forceApi.oauth.Authenticate(ctx)
	if err != nil {
		t.Fatalf("Unable to authenticate: %#v", err)
	}
	if err := forceApi.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// We shouldn't hit any errors creating a new force instance and manually passing in these oauth details now.
	newForceApi, err := CreateWithAccessToken(ctx, testVersion, testClientId, forceApi.oauth.AccessToken, forceApi.oauth.InstanceUrl, testLoginUri)
	if err != nil {
		t.Fatalf("Unable to create new force api instance using pre-defined oauth details: %#v", err)
	}
	if err := newForceApi.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// We should be able to make a basic query now with the newly created object (i.e. the oauth details should be correctly usable).
	_, err = newForceApi.DescribeSObject(ctx, &sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}
}
