package force

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	version      = "1.0.0"
	userAgent    = "go-force/" + version
	contentType  = "application/json"
	responseType = "application/json"
)

// Get issues a GET to the specified path with the given params and put the
// umarshalled (json) result in the third parameter
func (forceAPI *ForceAPI) Get(path string, params url.Values, out interface{}) error {
	return forceAPI.request("GET", path, params, nil, out)
}

// Post issues a POST to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceAPI *ForceAPI) Post(path string, params url.Values, payload, out interface{}) error {
	return forceAPI.request("POST", path, params, payload, out)
}

// Put issues a PUT to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceAPI *ForceAPI) Put(path string, params url.Values, payload, out interface{}) error {
	return forceAPI.request("PUT", path, params, payload, out)
}

// Patch issues a PATCH to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceAPI *ForceAPI) Patch(path string, params url.Values, payload, out interface{}) error {
	return forceAPI.request("PATCH", path, params, payload, out)
}

// Delete issues a DELETE to the specified path with the given payload
func (forceAPI *ForceAPI) Delete(path string, params url.Values) error {
	return forceAPI.request("DELETE", path, params, nil, nil)
}

func (forceAPI *ForceAPI) request(method, path string, params url.Values, payload, out interface{}) error {
	if err := forceAPI.oauth.Validate(); err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(forceAPI.oauth.InstanceUrl)
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	// Build body
	var body io.Reader
	if payload != nil {

		jsonBytes, err := json.Marshal(payload)
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
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", forceAPI.oauth.AccessToken))

	// Send
	forceAPI.traceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending %v request: %v", method, err)
	}
	defer resp.Body.Close()
	forceAPI.traceResponse(resp)

	// Sometimes the force API returns no body, we should catch this early
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response bytes: %v", err)
	}
	forceAPI.traceResponseBody(respBytes)

	// Attempt to parse response into out
	var objectUnmarshalErr error
	if out != nil {
		objectUnmarshalErr = json.Unmarshal(respBytes, out)
		if objectUnmarshalErr == nil {
			return nil
		}
	}

	// Attempt to parse response as a force.com api error before returning object unmarshal err
	apiErrors := APIErrors{}
	if marshalErr := json.Unmarshal(respBytes, &apiErrors); marshalErr == nil {
		if apiErrors.Validate() {
			// Check if error is oauth token expired
			if forceAPI.oauth.Expired(apiErrors) {
				// Reauthenticate then attempt query again
				oauthErr := forceAPI.oauth.Authenticate()
				if oauthErr != nil {
					return oauthErr
				}

				return forceAPI.request(method, path, params, payload, out)
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

func (forceAPI *ForceAPI) traceRequest(req *http.Request) {
	if forceAPI.logger != nil {
		forceAPI.trace("Request:", req, "%v")
	}
}

func (forceAPI *ForceAPI) traceResponse(resp *http.Response) {
	if forceAPI.logger != nil {
		forceAPI.trace("Response:", resp, "%v")
	}
}

func (forceAPI *ForceAPI) traceResponseBody(body []byte) {
	if forceAPI.logger != nil {
		forceAPI.trace("Response Body:", string(body), "%s")
	}
}
