package force

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nimajalali/go-force/forcejson"
)

const (
	version      = "1.0.0"
	userAgent    = "go-force/" + version
	contentType  = "application/json"
	responseType = "application/json"
)

func (forceApi *ForceApi) get(path string, payload url.Values, out interface{}) error {
	return forceApi.request("GET", path, payload, nil, out)
}

func (forceApi *ForceApi) post(path string, payload url.Values, body, out interface{}) error {
	return forceApi.request("POST", path, payload, body, out)
}

func (forceApi *ForceApi) patch(path string, payload url.Values, body, out interface{}) error {
	return forceApi.request("PATCH", path, payload, body, out)
}

func (forceApi *ForceApi) delete(path string, payload url.Values) error {
	return forceApi.request("DELETE", path, payload, nil, nil)
}

func (forceApi *ForceApi) request(method, path string, params url.Values, payload, out interface{}) error {
	if err := forceApi.oauth.Validate(); err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(forceApi.oauth.InstanceUrl)
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	// Build body
	var body io.Reader
	if payload != nil {

		jsonBytes, err := forcejson.Marshal(payload)
		if err != nil {
			return fmt.Errorf("Error marshaling encoded payload: %v", err)
		}

		body = bytes.NewReader(jsonBytes)
	}

	// Build Request
	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", responseType)
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", forceApi.oauth.AccessToken))

	// Send
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending %v request: %v", method, err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response bytes: %v", err)
	}

	// Attempt to parse response into out
	var objectUnmarshalErr error
	if out != nil {
		objectUnmarshalErr = forcejson.Unmarshal(respBytes, out)
		if objectUnmarshalErr == nil {
			return nil
		}
	}

	// Attempt to parse response as a force.com api error before returning object unmarshal err
	apiErrors := ApiErrors{}
	if marshalErr := forcejson.Unmarshal(respBytes, &apiErrors); marshalErr == nil {
		if apiErrors.Validate() {
			// Check if error is oauth token expired
			if forceApi.oauth.Expired(apiErrors) {
				// Reauthenticate then attempt query again
				oauthErr := forceApi.oauth.Authenticate()
				if oauthErr != nil {
					return oauthErr
				}

				return forceApi.request(method, path, params, payload, out)
			}

			return apiErrors
		}
	}

	if objectUnmarshalErr != nil {
		// Not a force.com api error. Just an unmarshalling error.
		return fmt.Errorf("Unable to unmarshal response to object: %v", objectUnmarshalErr)
	}

	// Sometimes no response is expected. For example delete and update. We still have to make sure an error wasn't returned.
	return nil
}
