// Package force provides bindings to the force.com REST API
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

var requestedObjMetadata = make(map[string]struct{})

// setRequestedObjectMetadata set the list of SObjects
// for which metadata have to be fetched.
func setRequestedObjectMetadata(list []string) {
	if list == nil {
		return
	}
	for _, obj := range list {
		requestedObjMetadata[obj] = struct{}{}
	}
}

// Create initialises a new SalesForce API client.
func Create(version, clientID, clientSecret, userName, password, securityToken,
	environment string, requiredObj []string) (*API, error) {
	oauth := &forceOauth{
		clientID:      clientID,
		clientSecret:  clientSecret,
		userName:      userName,
		password:      password,
		securityToken: securityToken,
		environment:   environment,
	}

	setRequestedObjectMetadata(requiredObj)

	forceAPI := &API{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
		openJobMap:             make(map[string]SJob),
	}

	// Init oauth
	err := forceAPI.oauth.Authenticate()
	if err != nil {
		return nil, err
	}

	// Init Api Resources
	err = forceAPI.getAPIResources()
	if err != nil {
		return nil, fmt.Errorf("Cannot get API resources: %s", err)
	}
	err = forceAPI.getAPISObjects()
	if err != nil {
		return nil, fmt.Errorf("Cannot get SObjects: %s", err)
	}
	err = forceAPI.getAPISObjectDescriptions()
	if err != nil {
		return nil, fmt.Errorf("Cannot get SObject descriptions: %s", err)
	}

	return forceAPI, nil
}

// CreateWithAccessToken initialises a new SalesForce API client with an access token.
func CreateWithAccessToken(version, clientID, accessToken, instanceURL string) (*API, error) {
	oauth := &forceOauth{
		clientID:    clientID,
		AccessToken: accessToken,
		InstanceURL: instanceURL,
	}

	forceAPI := &API{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
		openJobMap:             make(map[string]SJob),
	}

	// We need to check for oath correctness here, since we are not generating the token ourselves.
	if err := forceAPI.oauth.Validate(); err != nil {
		return nil, err
	}

	// Init Api Resources
	err := forceAPI.getAPIResources()
	if err != nil {
		return nil, err
	}
	err = forceAPI.getAPISObjects()
	if err != nil {
		return nil, err
	}

	return forceAPI, nil
}

// CreateWithRefreshToken initialises a new SalesForce API client with a refresh token.
func CreateWithRefreshToken(version, clientID, accessToken, instanceURL string) (*API, error) {
	oauth := &forceOauth{
		clientID:    clientID,
		AccessToken: accessToken,
		InstanceURL: instanceURL,
	}

	forceAPI := &API{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             version,
		oauth:                  oauth,
		openJobMap:             make(map[string]SJob),
	}

	// obtain access token
	if err := forceAPI.RefreshToken(); err != nil {
		return nil, err
	}

	// We need to check for oath correctness here, since we are not generating the token ourselves.
	if err := forceAPI.oauth.Validate(); err != nil {
		return nil, err
	}

	// Init Api Resources
	err := forceAPI.getAPIResources()
	if err != nil {
		return nil, err
	}
	err = forceAPI.getAPISObjects()
	if err != nil {
		return nil, err
	}

	return forceAPI, nil
}

// Used when running tests.
func createTest() *API {
	forceAPI, err := Create(testVersion, testClientID, testClientSecret, testUserName, testPassword, testSecurityToken, testEnvironment, nil)
	if err != nil {
		fmt.Printf("Unable to create API for test: %v", err)
		os.Exit(1)
	}

	return forceAPI
}

// APILogger describes an API logger.
type APILogger interface {
	Printf(format string, v ...interface{})
}

// TraceOn turns on logging for this API. After this is called, all
// requests, responses, and raw response bodies will be sent to the logger.
// If prefix is a non-empty string, it will be written to the front of all
// logged strings, which can aid in filtering log lines.
//
// Use TraceOn if you want to spy on the API requests and responses.
//
// Note that the base log.Logger type satisfies APILogger, but adapters
// can easily be written for other logging packages (e.g., the
// golang-sanctioned glog framework).
func (forceAPI *API) TraceOn(prefix string, logger APILogger) {
	forceAPI.logger = logger
	if prefix == "" {
		forceAPI.logPrefix = prefix
	} else {
		forceAPI.logPrefix = fmt.Sprintf("%s ", prefix)
	}
}

// TraceOff turns off tracing. It is idempotent.
func (forceAPI *API) TraceOff() {
	forceAPI.logger = nil
	forceAPI.logPrefix = ""
}

func (forceAPI *API) trace(name string, value interface{}, format string) {
	if forceAPI.logger != nil {
		logMsg := "%s%s " + format + "\n"
		forceAPI.logger.Printf(logMsg, forceAPI.logPrefix, name, value)
	}
}
