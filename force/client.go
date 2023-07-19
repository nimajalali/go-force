package force

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pwaterz/go-force/forcejson"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	version      = "1.0.0"
	userAgent    = "go-force/" + version
	contentType  = "application/json"
	responseType = "application/json"
)

type CompositeRequest struct {
	SubRequests []CompositeSubRequest `json:"compositeRequest"`
}

type CompositeSubRequest struct {
	Method       string          `json:"method"`
	Url          string          `json:"url"`
	ReferencedID string          `json:"referenceId"`
	HttpHeaders  http.Header     `json:"httpHeaders,omitempty"`
	Body         json.RawMessage `json:"body,omitempty"`
}

type CompositeResponse struct {
	Response []Subvalue `json:"compositeResponse"`
}

// Subvalue is the subresponses to the composite API.  Using the
// referende id, one will be able to match the response with the request.
type Subvalue struct {
	Body           json.RawMessage   `json:"body"`
	HTTPHeaders    map[string]string `json:"httpHeaders"`
	HTTPStatusCode int               `json:"httpStatusCode"`
	ReferenceID    string            `json:"referenceId"`
}

// Get issues a GET to the specified path with the given params and put the
// umarshalled (json) result in the third parameter
func (forceApi *ForceApi) Get(ctx context.Context, path string, params url.Values, out interface{}) error {
	return forceApi.request(ctx, "GET", path, params, nil, out, false)
}

// Post issues a POST to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Post(ctx context.Context, path string, params url.Values, payload, out interface{}) error {
	return forceApi.request(ctx, "POST", path, params, payload, out, false)
}

// Put issues a PUT to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Put(ctx context.Context, path string, params url.Values, payload, out interface{}) error {
	return forceApi.request(ctx, "PUT", path, params, payload, out, false)
}

// Patch issues a PATCH to the specified path with the given params and payload
// and put the unmarshalled (json) result in the third parameter
func (forceApi *ForceApi) Patch(ctx context.Context, path string, params url.Values, payload, out interface{}) error {
	return forceApi.request(ctx, "PATCH", path, params, payload, out, false)
}

// Delete issues a DELETE to the specified path with the given payload
func (forceApi *ForceApi) Delete(ctx context.Context, path string, params url.Values) error {
	return forceApi.request(ctx, "DELETE", path, params, nil, nil, false)
}

func (forceApi *ForceApi) request(ctx context.Context, method, path string, params url.Values, payload, out interface{}, retry bool) error {
	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(path)
	if len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}
	uriString := uri.String()
	compSubRequest := CompositeSubRequest{
		Method:       method,
		Url:          uriString,
		ReferencedID: "sf-request",
	}

	if payload != nil {
		jsonBytes, err := forcejson.Marshal(payload)
		if err != nil {
			return fmt.Errorf("Error marshaling encoded payload: %v", err)
		}
		compSubRequest.Body = jsonBytes
	}

	// Build Request
	compRequest := CompositeRequest{
		SubRequests: []CompositeSubRequest{compSubRequest},
	}

	jsonBytes, err := json.Marshal(compRequest)
	if err != nil {
		return fmt.Errorf("Error marshaling encoded payload: %v", err)
	}
	var sfURI bytes.Buffer
	sfURI.WriteString(forceApi.InstanceURL)
	sfURI.WriteString("/services/data/" + forceApi.apiVersion + "/composite")

	req, err := http.NewRequestWithContext(ctx, "POST", sfURI.String(), bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}
	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", responseType)

	// Set this for this request only if requested by caller (if this header is set, OwnerId of created case
	// will be set to the one we are passing in the request; if header not set, OwnerId is overwritten using SF rules)
	if forceApi.disableForceAutoAssign {
		fmt.Println("Disabling force auto assign")
		req.Header.Set("Sforce-Auto-Assign", "False")
		forceApi.SetDisableForceAutoAssign(false)
	}

	// Send
	forceApi.traceRequest(req)
	span, ctx := tracer.StartSpanFromContext(ctx, "Salesforce API Request", tracer.ResourceName(uriString))
	span.SetTag("url", path)
	span.SetTag("params", params)
	span.SetTag("http_method", method)
	span.SetTag("http_body", compSubRequest.Body)
	resp, err := forceApi.client.Do(req)
	if err != nil {
		returnErr := fmt.Errorf("Error sending %v request: %v", method, err)
		span.Finish(tracer.WithError(returnErr))
		return returnErr
	}
	span.Finish()
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
	if resp.StatusCode >= 300 {
		// Attempt to parse response as a force.com api error before returning object unmarshal err
		apiErrors := ApiErrors{}
		if marshalErr := forcejson.Unmarshal(respBytes, &apiErrors); marshalErr == nil {
			if apiErrors.Validate() {
				// Deal with expired salesforce tokens
				if apiErrors[0].ErrorCode == "INVALID_SESSION_ID" && !retry {
					forceApi.jwtMutex.Lock()
					forceApi.client = forceApi.jwtConfig.Client(ctx)
					forceApi.jwtMutex.Unlock()
					return forceApi.request(ctx, method, path, params, payload, out, true)
				}
				span.Finish(tracer.WithError(apiErrors))

				return apiErrors
			}
		}
	}

	forceApi.traceResponseBody(respBytes)
	compResponse := &CompositeResponse{}
	err = json.Unmarshal(respBytes, compResponse)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return fmt.Errorf("Error unmarshaling compsite response: %v", err)
	}
	if len(compResponse.Response) == 0 {
		return errors.New("sf response is empty")
	}

	if compResponse.Response[0].HTTPStatusCode >= 300 {
		// Attempt to parse response as a force.com api error before returning object unmarshal err
		apiErrors := ApiErrors{}
		if marshalErr := forcejson.Unmarshal(compResponse.Response[0].Body, &apiErrors); marshalErr == nil {
			if apiErrors.Validate() {
				span.Finish(tracer.WithError(apiErrors))
				return apiErrors
			}
		}
	}

	// Attempt to parse response into out
	var objectUnmarshalErr error
	if out != nil {
		objectUnmarshalErr = forcejson.Unmarshal(compResponse.Response[0].Body, out)
		if objectUnmarshalErr == nil {
			return nil
		}
	}

	if objectUnmarshalErr != nil {
		span.Finish(tracer.WithError(objectUnmarshalErr))
		// Not a force.com api error. Just an unmarshalling error.
		return fmt.Errorf("Unable to unmarshal response to object: %v", objectUnmarshalErr)
	}

	// Sometimes no response is expected. For example delete and update. We still have to make sure an error wasn't returned.
	return nil
}

func (forceApi *ForceApi) requestNonComposite(ctx context.Context, method, path string, params url.Values, payload, out interface{}, retry bool) error {
	// Build Uri
	var uri bytes.Buffer
	uri.WriteString(forceApi.InstanceURL)
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
	uriString := uri.String()
	req, err := http.NewRequestWithContext(ctx, method, uriString, body)
	if err != nil {
		return fmt.Errorf("Error creating %v request: %v", method, err)
	}

	// Add Headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", responseType)

	// Set this for this request only if requested by caller (if this header is set, OwnerId of created case
	// will be set to the one we are passing in the request; if header not set, OwnerId is overwritten using SF rules)
	if forceApi.disableForceAutoAssign {
		fmt.Println("Disabling force auto assign")
		req.Header.Set("Sforce-Auto-Assign", "False")
		forceApi.SetDisableForceAutoAssign(false)
	}

	// Send
	forceApi.traceRequest(req)
	span, ctx := tracer.StartSpanFromContext(ctx, "Salesforce API Request", tracer.ResourceName(uriString))
	span.SetTag("url", uriString)
	span.SetTag("http_method", method)
	resp, err := forceApi.client.Do(req)
	if err != nil {
		returnErr := fmt.Errorf("Error sending %v request: %v", method, err)
		span.Finish(tracer.WithError(returnErr))
		return returnErr
	}
	span.Finish()
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
			// Deal with expired salesforce tokens
			if apiErrors[0].ErrorCode == "INVALID_SESSION_ID" && !retry {
				forceApi.jwtMutex.Lock()
				forceApi.client = forceApi.jwtConfig.Client(ctx)
				forceApi.jwtMutex.Unlock()
				return forceApi.request(ctx, method, path, params, payload, out, true)
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
