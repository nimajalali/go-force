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
	contentType      = "application/json; charset=UTF-8"
	gzipEncodingType = "gzip"
	responseType     = "application/json"
	version          = "1.0.0"
	userAgent        = "go-force/" + version
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
			fmt.Println("SSSSSSSSSSSSSSSSSSSSSSSSSSSTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT: ", payload.(string))
			body = bytes.NewReader([]byte(payload.(string)))
		default:
			jsonBytes, err := forcejson.Marshal(payload)
			//jsonBytes, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("Error marshalling encoded payload: %v", err)
			}
			fmt.Println("----------------SSSSSSSSSSSSSSSSSSSSSSSSSSSTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT: ", string(jsonBytes))
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
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", responseType)
	req.Header.Set("Accept-Encoding", gzipEncodingType)
	req.Header.Set("Authorization", fmt.Sprintf("%v %v", "Bearer", forceAPI.oauth.AccessToken))
	req.Header.Set("X-SFDC-Session", forceAPI.oauth.AccessToken)

	// Send
	fmt.Printf("-------RRRRRRRRRRRR: %+v\n", req)

	forceAPI.traceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending %v request: %v", method, err)
	}

	forceAPI.traceResponse(resp)

	// Sometimes the force API returns no body, we should catch this early
	if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("NNNNNNNNNNHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH:: %+v\n", resp)
		return nil
	}

	fmt.Printf("CCCCCCCCCCCCCCCCCCCC: %+v\n", resp)

	var respBytes []byte
	// gzip decoding
	if resp.Header.Get("Content-Encoding") == gzipEncodingType {

		fmt.Println("GGGGGGGGGGGGGZIIIIIIIIIIIPPPP")
		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("Cannot decode gzip response: %s", err)
		}
		respBytes, err = ioutil.ReadAll(zr)
		if err != nil {
			return fmt.Errorf("Cannot read decoded gzip response: %s", err)
		}
		err = zr.Close()
		if err != nil {
			return fmt.Errorf("Cannot close gzip reader: %s", err)
		}
	} else {
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading response bytes: %v", err)
		}
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Bad response: %s", string(respBytes))
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Cannot close request: %s", err)
	}

	forceAPI.traceResponseBody(respBytes)

	fmt.Printf("----VVVVVVVVVVV: %+v\n\n", string(respBytes))

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
	if marshalErr := forcejson.Unmarshal(respBytes, &apiErrors); marshalErr == nil {
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
