// A Go package that provides bindings to the force.com REST API
//
// See http://www.salesforce.com/us/developer/docs/api_rest/
package force

const (
	grantType = "password"

	testVersion       = ""
	testClientId      = ""
	testClientSecret  = ""
	testUserName      = ""
	testPassword      = ""
	testSecurityToken = ""
	testEnvironment   = "sandbox"
)

// Basic information needed to connect to the Force.com REST API.
var apiVersion string
var oauth *ForceOauth

func Init(version, clientId, clientSecret, username, password, securityToken, environment string) error {
	apiVersion = version

	var err error
	oauth, err = authenticate(grantType, clientId, clientSecret, username, password, securityToken, environment)
	if err != nil {
		return err
	}

	return nil
}

// Used when running tests.
func initTest() error {
	apiVersion = testVersion

	var err error
	oauth, err = authenticate(grantType, testClientId, testClientSecret, testUserName, testPassword, testSecurityToken, testEnvironment)
	if err != nil {
		return err
	}

	return nil
}
