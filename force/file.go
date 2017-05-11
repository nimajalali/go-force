package force

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/simplesurance/go-force/forcejson"
)

// FileUpload creates a new file upload http request with optional extra params
func (forceAPI *API) FileUpload(j *SJob, params map[string]string, paramName string, file *os.File) (*SBatch, error) {

	fi, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Cannot get request file info: %s", err)
	}

	URI := fmt.Sprintf("%s/%s/batch", j.forceAPI.oauth.InstanceURL, j.BaseURI)
	req, err := http.NewRequest("POST", URI, file)

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", responseType)
	req.Header.Set("Accept-Encoding", gzipEncodingType)
	req.Header.Set("X-SFDC-Session", j.forceAPI.oauth.AccessToken)
	req.Header.Set("Content-Type", zipJSONContentType)
	req.Header.Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	forceAPI.traceRequest(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Cannot send file upload request: %s", err)
	}
	forceAPI.traceResponse(resp)

	batch := &SBatch{}

	// Sometimes (for updates) the force API returns no body, we should catch this early
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var respBytes []byte
	// gzip decoding
	if resp.Header.Get("Content-Encoding") == gzipEncodingType {
		respBytes, err = gzipDecode(resp.Body)
	} else {
		respBytes, err = ioutil.ReadAll(resp.Body)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Cannot close request body: %s", err)
	}

	forceAPI.traceResponseBody(respBytes)

	// Attempt to parse response into out
	var objectUnmarshalErr error
	objectUnmarshalErr = json.Unmarshal(respBytes, batch)
	if objectUnmarshalErr == nil {
		return batch, nil
	}

	// Attempt to parse response as a force.com api error before returning object unmarshal err
	apiErrors := APIErrors{}
	if marshalErr := forcejson.Unmarshal(respBytes, &apiErrors); marshalErr == nil {
		if apiErrors.Validate() {
			// Check if error is oauth token expired
			if forceAPI.oauth.Expired(apiErrors) {
				forceAPI.logger.Printf("ForceApi session token has expired")

				// Reauthenticate then attempt query again
				forceAPI.logger.Printf("Trying to get a new session token")
				oauthErr := forceAPI.oauth.Authenticate()
				if oauthErr != nil {
					forceAPI.logger.Printf("Failed to get a new session token: %v", oauthErr)
					return nil, oauthErr
				}
				return forceAPI.FileUpload(j, params, paramName, file)
			}

			return nil, apiErrors
		}
	}

	if objectUnmarshalErr != nil {
		// Not a force.com api error. Just an unmarshalling error.
		return nil, fmt.Errorf("Unable to unmarshal response to object: %v", objectUnmarshalErr)
	}

	// Sometimes no response is expected. For example delete and update. We still have to make sure an error wasn't returned.
	return batch, nil
}
