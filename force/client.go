package force

import (
	"bytes"

	"github.com/nimajalali/go-force/force/encoding"
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

func request(method, path string, payload url.Values, body, out interface{}) error {
	if err := oauth.Validate(); err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(oauth.InstanceUrl)
	uri.WriteString(path)
	if payload != nil && len(payload) != 0 {
		uri.WriteString("?")
		uri.WriteString(payload.Encode())
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

	// Attempt to parse response as a salesforce api error
	apiError := ApiError{}
	if err := json.Unmarshal(respBytes, &apiError); err == nil {
		// Check if api error is valid
		if len(apiError.ErrorCode) != 0 {
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
