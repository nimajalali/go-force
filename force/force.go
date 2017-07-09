// A Go package that provides bindings to the force.com REST API
//
// See http://www.salesforce.com/us/developer/docs/api_rest/
package force

import (
	"fmt"
	"os"
)

const (
	testVersion       = "v36.0"
	testClientID      = "3MVG9A2kN3Bn17hs8MIaQx1voVGy662rXlC37svtmLmt6wO_iik8Hnk3DlcYjKRvzVNGWLFlGRH1ryHwS217h"
	testClientSecret  = "4165772184959202901"
	testUserName      = "go-force@jalali.net"
	testPassword      = "golangrocks3"
	testSecurityToken = "kAlicVmti9nWRKRiWG3Zvqtte"
	testEnvironment   = "production"
)

func Create(version, clientID, clientSecret, userName, password, securityToken,
	environment string) (*ForceAPI, error) {
	oauth := &forceOauth{
		clientID:      clientID,
		clientSecret:  clientSecret,
		userName:      userName,
		password:      password,
		securityToken: securityToken,
		environment:   environment,
	}

	ForceAPI := &ForceAPI{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	// Init oauth
	err := ForceAPI.oauth.Authenticate()
	if err != nil {
		return nil, err
	}

	// Init API Resources
	err = ForceAPI.getResources()
	if err != nil {
		return nil, err
	}
	err = ForceAPI.getSObjects()
	if err != nil {
		return nil, err
	}

	return ForceAPI, nil
}

func CreateWithAccessToken(version, clientID, accessToken, instanceUrl string) (*ForceAPI, error) {
	oauth := &forceOauth{
		clientID:    clientID,
		AccessToken: accessToken,
		InstanceUrl: instanceUrl,
	}

	ForceAPI := &ForceAPI{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	// We need to check for oath correctness here, since we are not generating the token ourselves.
	if err := ForceAPI.oauth.Validate(); err != nil {
		return nil, err
	}

	// Init API Resources
	err := ForceAPI.getResources()
	if err != nil {
		return nil, err
	}
	err = ForceAPI.getSObjects()
	if err != nil {
		return nil, err
	}

	return ForceAPI, nil
}

func CreateWithRefreshToken(version, clientID, clientSecret, refreshToken, instanceUrl string) (*ForceAPI, error) {
	oauth := &forceOauth{
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		InstanceUrl:  instanceUrl,
	}

	ForceAPI := &ForceAPI{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
	}

	// obtain access token
	if err := ForceAPI.RefreshToken(); err != nil {
		return nil, err
	}

	// We need to check for oath correctness here, since we are not generating the token ourselves.
	if err := ForceAPI.oauth.Validate(); err != nil {
		return nil, err
	}

	// Init API Resources
	err := ForceAPI.getResources()
	if err != nil {
		return nil, err
	}
	err = ForceAPI.getSObjects()
	if err != nil {
		return nil, err
	}

	return ForceAPI, nil
}

// Used when running tests.
func createTest() *ForceAPI {
	ForceAPI, err := Create(testVersion, testClientID, testClientSecret, testUserName, testPassword, testSecurityToken, testEnvironment)
	if err != nil {
		fmt.Printf("Unable to create ForceAPI for test: %v", err)
		os.Exit(1)
	}

	return ForceAPI
}

type ForceAPILogger interface {
	Printf(format string, v ...interface{})
}

// TraceOn turns on logging for this ForceAPI. After this is called, all
// requests, responses, and raw response bodies will be sent to the logger.
// If prefix is a non-empty string, it will be written to the front of all
// logged strings, which can aid in filtering log lines.
//
// Use TraceOn if you want to spy on the ForceAPI requests and responses.
//
// Note that the base log.Logger type satisfies ForceAPILogger, but adapters
// can easily be written for other logging packages (e.g., the
// golang-sanctioned glog framework).
func (ForceAPI *ForceAPI) TraceOn(prefix string, logger ForceAPILogger) {
	ForceAPI.logger = logger
	if prefix == "" {
		ForceAPI.logPrefix = prefix
	} else {
		ForceAPI.logPrefix = fmt.Sprintf("%s ", prefix)
	}
}

// TraceOff turns off tracing. It is idempotent.
func (ForceAPI *ForceAPI) TraceOff() {
	ForceAPI.logger = nil
	ForceAPI.logPrefix = ""
}

func (ForceAPI *ForceAPI) trace(name string, value interface{}, format string) {
	if ForceAPI.logger != nil {
		logMsg := "%s%s " + format + "\n"
		ForceAPI.logger.Printf(logMsg, ForceAPI.logPrefix, name, value)
	}
}
