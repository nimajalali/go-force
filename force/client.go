package force

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/simplesurance/go-force/forcejson"
)

const (
	jsonContentType    = "application/json; charset=UTF-8"
	zipJSONContentType = "zip/json"
	gzipEncodingType   = "gzip"
	responseType       = "application/json"
	version            = "1.0.0"
	userAgent          = "go-force/" + version
)

// Get issues a GET to the specified path with the given params and put the
// umarshalled (json) result in the third parameter
func (forceAPI *API) Get(path string, params url.Values, out interface{}) error {
	return forceAPI.request("GET", path, params, nil, out)
}

// Post issues a POST to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceAPI *API) Post(path string, params url.Values, payload, out interface{}) error {
	return forceAPI.request("POST", path, params, payload, out)
}

// Put issues a PUT to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceAPI *API) Put(path string, params url.Values, payload, out interface{}) error {
	return forceAPI.request("PUT", path, params, payload, out)
}

// Patch issues a PATCH to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceAPI *API) Patch(path string, params url.Values, payload, out interface{}) error {
	return forceAPI.request("PATCH", path, params, payload, out)
}

// Delete issues a DELETE to the specified path with the given payload
func (forceAPI *API) Delete(path string, params url.Values) error {
	return forceAPI.request("DELETE", path, params, nil, nil)
}

func gzipDecode(body io.Reader) ([]byte, error) {
	zr, err := gzip.NewReader(body)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode gzip response: %s", err)
	}
	buf, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("Cannot read decoded gzip response: %s", err)
	}
	err = zr.Close()
	if err != nil {
		return nil, fmt.Errorf("Cannot close gzip reader: %s", err)
	}
	return buf, nil
}

func (forceAPI *API) request(method, path string, params url.Values, payload, out interface{}) error {
	if err := forceAPI.oauth.Validate(); err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(forceAPI.oauth.InstanceURL)
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	// Build body
	var body io.Reader
	if payload != nil {

		switch payload.(type) {
		case string:
			body = bytes.NewReader([]byte(payload.(string)))
		default:
			jsonBytes, err := forcejson.Marshal(payload)
			if err != nil {
				return fmt.Errorf("Error marshalling encoded payload: %v", err)
			}
			body = bytes.NewReader(jsonBytes)
		}
	}

	// Build Request
	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Accept", responseType)
	req.Header.Set("Accept-Encoding", gzipEncodingType)
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", forceAPI.oauth.AccessToken))
	req.Header.Set("X-SFDC-Session", forceAPI.oauth.AccessToken)

	// Read the content
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}
	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("Error reading %v request: %v", method, err)
	}

	forceAPI.traceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending %v request: %v", method, err)
	}

	return forceAPI.readResponse(resp, method, path, params, payload, out)
}

func (forceAPI *API) readResponse(resp *http.Response, method, path string, params url.Values, payload, out interface{}) error {
	forceAPI.traceResponse(resp)

	// Sometimes (for updates) the force API returns no body, we should catch this early
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var body []byte
	var err error
	// gzip decoding
	if resp.Header.Get("Content-Encoding") == gzipEncodingType {
		body, err = gzipDecode(resp.Body)
	} else {
		body, err = ioutil.ReadAll(resp.Body)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Cannot close request body: %s", err)
	}

	err = forceAPI.processResponse(body, method, path, params, payload, out)
	if err != nil {
		return fmt.Errorf("Cannot process response: %s", err)
	}

	// Sometimes no response is expected. For example delete and update. We still have to make sure an error wasn't returned.
	return nil
}

func (forceAPI *API) processResponse(body []byte, method, path string, params url.Values, payload, out interface{}) error {
	forceAPI.traceResponseBody(body)

	// Attempt to parse response into out
	var objectUnmarshalErr error
	if out != nil {
		objectUnmarshalErr = json.Unmarshal(body, out)
		if objectUnmarshalErr == nil {
			return nil
		}
	}

	// Attempt to parse response as a force.com api error before returning object unmarshal err
	apiErrors := APIErrors{}
	if marshalErr := forcejson.Unmarshal(body, &apiErrors); marshalErr == nil {
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
	return nil
}

func (forceAPI *API) traceRequest(req *http.Request) {
	if forceAPI.logger != nil {
		forceAPI.trace("Request:", req, "%v")
	}
}

func (forceAPI *API) traceResponse(resp *http.Response) {
	if forceAPI.logger != nil {
		forceAPI.trace("Response:", resp, "%v")
	}
}

func (forceAPI *API) traceResponseBody(body []byte) {
	if forceAPI.logger != nil {
		forceAPI.trace("Response Body:", string(body), "%s")
	}
}
