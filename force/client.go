package force

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"github.com/nimajalali/go-force/forcejson"
)

const (
	version      	   = "1.0.0"
	userAgent          = "go-force/" + version
	refreshContentType = "application/x-www-form-urlencoded"
	contentType        = "application/json"
	responseType       = "application/json"
)

// Get issues a GET to the specified path with the given params and put the
// umarshalled (json) result in the third parameter
func (forceApi *ForceApi) Get(path string, params url.Values, out interface{}) error {
	return forceApi.request("GET", path, params, nil, out)
}

// Post issues a POST to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Post(path string, params url.Values, payload, out interface{}) error {
	_, refresh_token := payload.(map[string]string)["refresh_token"]

	if refresh_token {
		return forceApi.request_access_token(path, params, payload, out)
	}

	return forceApi.request("POST", path, params, payload, out)
}

// Put issues a PUT to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Put(path string, params url.Values, payload, out interface{}) error {
	return forceApi.request("PUT", path, params, payload, out)
}

// Patch issues a PATCH to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Patch(path string, params url.Values, payload, out interface{}) error {
	return forceApi.request("PATCH", path, params, payload, out)
}

// Delete issues a DELETE to the specified path with the given payload
func (forceApi *ForceApi) Delete(path string, params url.Values) error {
	return forceApi.request("DELETE", path, params, nil, nil)
}

func (forceApi *ForceApi) request_access_token(path string, params url.Values, payload, out interface{}) error {
	mPayload, _ := payload.(map[string]string)

	var uri bytes.Buffer
	uri.WriteString(forceApi.oauth.InstanceUrl)
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	url_param := `grant_type=refresh_token&`
	url_param += `client_id=%s&`
	url_param += `client_secret=%s&`
	url_param += `refresh_token=%s`
	url_param = fmt.Sprintf(
		url_param,
		mPayload["client_id"],
		mPayload["client_secret"],
		mPayload["refresh_token"],
	)

	req, _ := http.NewRequest("POST", uri.String(), strings.NewReader(url_param))

	req.Header.Add("Content-Type", refreshContentType)
	req.Header.Add("Cache-Control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(out)

	return nil
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
	req.Header.Set("Accept", responseType)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", forceApi.oauth.AccessToken))

	// Send
	forceApi.traceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending %v request: %v", method, err)
	}
	defer resp.Body.Close()
	forceApi.traceResponse(resp)

	// Sometimes the force API returns no body, we should catch this early
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response bytes: %v", err)
	}
	forceApi.traceResponseBody(respBytes)

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

func (forceApi *ForceApi) traceRequest(req *http.Request) {
	if forceApi.logger != nil {
		forceApi.trace("Request:", req, "%v")
	}
}

func (forceApi *ForceApi) traceResponse(resp *http.Response) {
	if forceApi.logger != nil {
		forceApi.trace("Response:", resp, "%v")
	}
}

func (forceApi *ForceApi) traceResponseBody(body []byte) {
	if forceApi.logger != nil {
		forceApi.trace("Response Body:", string(body), "%s")
	}
}

