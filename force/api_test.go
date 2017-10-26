package force

import "testing"

func TestHasAccess(t *testing.T) {

	oauth := &forceOauth{
		clientId:      testClientId,
		clientSecret:  testClientSecret,
		userName:      testUserName,
		password:      testPassword,
		securityToken: testSecurityToken,
		environment:   testEnvironment,
	}

	forceApi := &ForceApi{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	apiName := "AbadCompliance___c"

	forceApi.apiSObjects[apiName] = &SObjectMetaData{}
	validObjects := []string{apiName}
	if value := forceApi.HasAccess(validObjects); !value {
		t.Error("expected to return true, but got false")
	}
	invalidObjects := []string{"Alien"}
	if value := forceApi.HasAccess(invalidObjects); value {
		t.Error("expected to return false, but got true")
	}
}
