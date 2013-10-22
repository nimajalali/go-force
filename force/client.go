package force

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"go-force/force/encoding"
)

const (
	version      = "0.0.1"
	userAgent    = "go-force/" + version
	contentType  = "application/json"
	responseType = "application/json"
)

func Get(path string, payload url.Values, out interface{}) error {
	return request("GET", path, payload, nil, out)
}

func Post(path string, payload url.Values, body, out interface{}) error {
	return request("POST", path, payload, body, out)
}

func Patch(path string, payload url.Values, body, out interface{}) error {
	return request("PATCH", path, payload, body, out)
}

func Delete(path string, payload url.Values, body, out interface{}) error {
	return request("DELETE", path, payload, body, out)
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

	// Add Auth
	req.SetBasicAuth("Bearer", oauth.AccessToken)

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
	apiError := &ApiError{}
	if err := json.Unmarshal(respBytes, apiError); err == nil {
		// Check if api error is valid
		if len(apiError.Fields) != 0 || len(apiError.Message) != 0 || len(apiError.ErrorCode) != 0 || len(apiError.ErrorName) != 0 || len(apiError.ErrorDescription) != 0 {
			return apiError
		}
	}

	// Attempt to parse response into out
	if out != nil {
		if err := json.Unmarshal(respBytes, out); err != nil {
			return fmt.Errorf("Unable to unmarshal response: %v", err)
		}
	}

	return nil
}
