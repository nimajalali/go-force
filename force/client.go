package force

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nimajalali/go-force/encoding"
)

const (
	version      = "0.0.1"
	userAgent    = "go-force/" + version
	contentType  = "application/json"
	responseType = "application/json"
)

func get(path string, payload url.Values, out interface{}) error {
	return request("GET", path, payload, nil, out)
}

func post(path string, payload url.Values, body, out interface{}) error {
	return request("POST", path, payload, body, out)
}

func patch(path string, payload url.Values, body, out interface{}) error {
	return request("PATCH", path, payload, body, out)
}

func delete(path string, payload url.Values) error {
	return request("DELETE", path, payload, nil, nil)
}

func request(method, path string, params url.Values, payload, out interface{}) error {
	if err := oauth.Validate(); err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(oauth.InstanceUrl)
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	// Build body
	var body io.Reader
	if payload != nil {
		encodedPayload, err := encoding.Encode(payload)
		if err != nil {
			return fmt.Errorf("Error encoding payload: %v", err)
		}

		jsonBytes, err := json.Marshal(encodedPayload)
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
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", oauth.AccessToken))

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

	// Attempt to parse response as a force.com api error
	apiErrors := ApiErrors{}
	err = json.Unmarshal(respBytes, &apiErrors)
	if err == nil {
		if apiErrors.Validate() {
			// Check if error is oauth token expired
			if oauth.Expired(apiErrors) {
				// Reauthenticate then attempt query again
				oauthErr := oauth.Authenticate()
				if oauthErr != nil {
					return oauthErr
				}

				return request(method, path, params, payload, out)
			}

			return apiErrors
		}
	}

	// Attempt to parse response into out
	if out != nil {
		// First parse json to map[string]interface{}
		respMap := make(map[string]interface{})
		if err := json.Unmarshal(respBytes, &respMap); err != nil {
			return fmt.Errorf("Unable to unmarshal response to map[string]interface{}: %v", err)
		}

		// Go-Force Decode
		if err := encoding.Decode(out, respMap); err != nil {
			return fmt.Errorf("Unable to dedode response to object: %v", err)
		}
	}

	return nil
}
