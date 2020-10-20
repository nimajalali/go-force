package force

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/EverlongProject/go-force/forcejson"
)

const (
	version      = "1.0.0"
	userAgent    = "go-force/" + version
	contentType  = "application/json"
	responseType = "application/json"
)

// Get issues a GET to the specified path with the given params and put the
// umarshalled (json) result in the third parameter
func (forceApi *ForceApi) Get(ctx context.Context, path string, params url.Values, out interface{}) error {
	return forceApi.request(ctx, "GET", path, params, nil, out)
}

// Post issues a POST to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Post(ctx context.Context, path string, params url.Values, payload, out interface{}) error {
	return forceApi.request(ctx, "POST", path, params, payload, out)
}

// Put issues a PUT to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Put(ctx context.Context, path string, params url.Values, payload, out interface{}) error {
	return forceApi.request(ctx, "PUT", path, params, payload, out)
}

// Patch issues a PATCH to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Patch(ctx context.Context, path string, params url.Values, payload, out interface{}) error {
	return forceApi.request(ctx, "PATCH", path, params, payload, out)
}

// Delete issues a DELETE to the specified path with the given payload
func (forceApi *ForceApi) Delete(ctx context.Context, path string, params url.Values) error {
	return forceApi.request(ctx, "DELETE", path, params, nil, nil)
}

func (forceApi *ForceApi) request(ctx context.Context, method, path string, params url.Values, payload, out interface{}) error {
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
	req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
	if err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", responseType)
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", forceApi.oauth.AccessToken))

	// Set this for this request only if requested by caller (if this header is set, OwnerId of created case
	// will be set to the one we are passing in the request; if header not set, OwnerId is overwritten using SF rules)
	if forceApi.disableForceAutoAssign {
		fmt.Println("Disabling force auto assign")
		req.Header.Set("Sforce-Auto-Assign", "False")
		forceApi.SetDisableForceAutoAssign(false)
	}

	// Send
	forceApi.traceRequest(req)
	resp, err := forceApi.client.Do(req)
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
				oauthErr := forceApi.oauth.Authenticate(ctx)
				if oauthErr != nil {
					return oauthErr
				}

				return forceApi.request(ctx, method, path, params, payload, out)
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
