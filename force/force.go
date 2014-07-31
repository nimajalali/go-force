// A Go package that provides bindings to the force.com REST API
//
// See http://www.salesforce.com/us/developer/docs/api_rest/
package force

import (
	"fmt"
	"os"
)

const (
	testVersion       = "v29.0"
	testClientId      = "3MVG9A2kN3Bn17hs8MIaQx1voVGy662rXlC37svtmLmt6wO_iik8Hnk3DlcYjKRvzVNGWLFlGRH1ryHwS217h"
	testClientSecret  = "4165772184959202901"
	testUserName      = "go-force@jalali.net"
	testPassword      = "golangrocks3"
	testSecurityToken = "kAlicVmti9nWRKRiWG3Zvqtte"
	testEnvironment   = "production"
)

// Basic information needed to connect to the Force.com REST API.
var apiVersion string
var oauth *forceOauth

func Init(version, clientId, clientSecret, userName, password, securityToken, environment string) error {
	apiVersion = version

	oauth = &forceOauth{
		clientId:      clientId,
		clientSecret:  clientSecret,
		userName:      userName,
		password:      password,
		securityToken: securityToken,
		environment:   environment,
	}

	// Init oauth
	err := oauth.Authenticate()
	if err != nil {
		return err
	}

	// Init Api Resources
	err = getApiResources()
	if err != nil {
		return err
	}
	err = getApiSObjects()
	if err != nil {
		return err
	}

	return nil
}

// Used when running tests.
func initTest() {
	// initTest is called multiple times throughout testing, it only needs to be run once.
	if len(apiVersion) == 0 {
		apiVersion = testVersion

		oauth = &forceOauth{
			clientId:      testClientId,
			clientSecret:  testClientSecret,
			userName:      testUserName,
			password:      testPassword,
			securityToken: testSecurityToken,
			environment:   testEnvironment,
		}

		// Init oauth
		err := oauth.Authenticate()
		if err != nil {
			fmt.Printf("Unable to authenticate for test: %v", err)
			os.Exit(1)
		}

		// Init Api Resources
		err = getApiResources()
		if err != nil {
			fmt.Printf("Unable to retrieve api resources for test: %v", err)
			os.Exit(1)
		}
		err = getApiSObjects()
		if err != nil {
			fmt.Printf("Unable to retrieve api sobjects for test: %v", err)
			os.Exit(1)
		}
	}
}
