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

func Create(version, clientId, clientSecret, userName, password, securityToken,
	environment string) (*ForceApi, error) {
	oauth := &forceOauth{
		clientId:      clientId,
		clientSecret:  clientSecret,
		userName:      userName,
		password:      password,
		securityToken: securityToken,
		environment:   environment,
	}

	forceApi := &ForceApi{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	// Init oauth
	err := forceApi.oauth.Authenticate()
	if err != nil {
		return nil, err
	}

	// Init Api Resources
	err = forceApi.getApiResources()
	if err != nil {
		return nil, err
	}
	err = forceApi.getApiSObjects()
	if err != nil {
		return nil, err
	}

	return forceApi, nil
}

func CreateWithAccessToken(version, clientId, accessToken, instanceUrl string) (*ForceApi, error) {
	oauth := &forceOauth{
		clientId:    clientId,
		AccessToken: accessToken,
		InstanceUrl: instanceUrl,
	}

	forceApi := &ForceApi{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	// We need to check for oath correctness here, since we are not generating the token ourselves.
	if err := forceApi.oauth.Validate(); err != nil {
		return nil, err
	}

	// Init Api Resources
	err := forceApi.getApiResources()
	if err != nil {
		return nil, err
	}
	err = forceApi.getApiSObjects()
	if err != nil {
		return nil, err
	}

	return forceApi, nil
}

// Used when running tests.
func createTest() *ForceApi {
	forceApi, err := Create(testVersion, testClientId, testClientSecret, testUserName, testPassword, testSecurityToken, testEnvironment)
	if err != nil {
		fmt.Printf("Unable to create ForceApi for test: %v", err)
		os.Exit(1)
	}

	return forceApi
}

type ForceApiLogger interface {
	Printf(format string, v ...interface{})
}

// TraceOn turns on logging for this ForceApi. After this is called, all
// requests, responses, and raw response bodies will be sent to the logger.
// If prefix is a non-empty string, it will be written to the front of all
// logged strings, which can aid in filtering log lines.
//
// Use TraceOn if you want to spy on the ForceApi requests and responses.
//
// Note that the base log.Logger type satisfies ForceApiLogger, but adapters
// can easily be written for other logging packages (e.g., the
// golang-sanctioned glog framework).
func (forceApi *ForceApi) TraceOn(prefix string, logger ForceApiLogger) {
	forceApi.logger = logger
	if prefix == "" {
		forceApi.logPrefix = prefix
	} else {
		forceApi.logPrefix = fmt.Sprintf("%s ", prefix)
	}
}

// TraceOff turns off tracing. It is idempotent.
func (forceApi *ForceApi) TraceOff() {
	forceApi.logger = nil
	forceApi.logPrefix = ""
}

func (forceApi *ForceApi) trace(name string, value interface{}, format string) {
	if forceApi.logger != nil {
		logMsg := "%s%s " + format + "\n"
		forceApi.logger.Printf(logMsg, forceApi.logPrefix, name, value)
	}
}
