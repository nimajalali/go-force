// A Go package that provides bindings to the force.com REST API
//
// See http://www.salesforce.com/us/developer/docs/api_rest/
package force

import (
	"github.com/nimajalali/go-force/force/oauth"
)

const (
	testApiVersion = ""

	testClientId      = ""
	testClientSecret  = ""
	testUserName      = ""
	testPassword      = ""
	testSecurityToken = ""
	testEnvironment   = "sandbox"

	grantType = "password"
)

// Basic information needed to connect to the Force.com REST API.
var version string
var oauth oauth.SFOauth

func Init(apiVersion, clientId, clientSecret, username, password, securityToken, environment string) error {
	version = apiVersion

	oauth, err := oauth.Authenticate(grantType, clientId, clientSecret, username, password, securityToken, environment)
	if err != nil {
		return err
	}

	return nil
}

// Used when running tests.
func initTest() error {
	version = testApiVersion

	oauth, err := oauth.Authenticate(grantType, testClientId, testClientSecret, testUserName, testPassword, testSecurityToken, testEnvironment)
	if err != nil {
		return err
	}

	return nil
}
