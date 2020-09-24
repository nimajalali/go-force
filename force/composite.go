package force

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

//CompositeRequests is a struct that contains all requests to be executed
//AllOrNone reverts all requests if a single request fails
type CompositeRequests struct {
	AllOrNone        bool               `force:"allOrNone,omitempty"`
	CompositeRequest []CompositeRequest `force:"compositeRequest,omitempty"`
}

//CompositeRequest defines a single request
type CompositeRequest struct {
	Method      string            `force:"method,omitempty"`
	URL         string            `force:"url,omitempty"`
	HTTPHeaders map[string]string `force:"httpHeaders,omitempty"`
	Body        interface{}       `force:"body,omitempty"`
	ReferenceID string            `force:"referenceId,omitempty"`
}

//CompositeResponses contains the responses to each request made
type CompositeResponses struct {
	CompositeResponse []CompositeResponse `force:"compositeResponse,omitempty"`
}

//CompositeResponse describes the response to a single request
type CompositeResponse struct {
	Body           json.RawMessage   `force:"body,omitempty"`
	HTTPHeaders    map[string]string `force:"httpHeaders,omitempty"`
	HTTPStatusCode int               `force:"httpStatusCode,omitempty"`
	ReferenceID    string            `force:"referenceId,omitempty"`
}

//PostCompositeRequests performs a composite request
//Salesforce API used must be v39.0 or greater
func (forceApi *ForceApi) PostCompositeRequests(in *CompositeRequests) (resp *CompositeResponses, err error) {
	uri := forceApi.apiResources[compositeKey]
	resp = &CompositeResponses{}

	err = forceApi.Post(uri, nil, in, resp)
	if err != nil {
		log.Debugf("Error returned from: forceApi.Post in PostCompositeRequests() func: %v\n", err)
	}

	return
}

//Add adds a request object to the CompositeRequests struct so that when it is executed, all containing requests are processed
func (requests *CompositeRequests) Add(request *CompositeRequest) {
	requests.CompositeRequest = append(requests.CompositeRequest, *request)
	return
}

//CompositeQuery generates a CompositeRequest in the format of a Query
func (forceApi *ForceApi) CompositeQuery(query string, referenceID string) *CompositeRequest {
	path := forceApi.apiResources[queryKey]

	params := url.Values{
		"q": {query},
	}

	var uri bytes.Buffer
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	return &CompositeRequest{
		Method:      http.MethodGet,
		URL:         uri.String(),
		ReferenceID: referenceID,
	}
}

//CompositeGetSObject generates a CompositeRequest in the format of a GetSObject
//@Params (obj SObject) only needs an empty object to define SObject type
func (forceApi *ForceApi) CompositeGetSObject(id string, obj SObject, fields []string, referenceID string) *CompositeRequest {

	path := strings.Replace(forceApi.apiSObjects[obj.ApiName()].URLs[rowTemplateKey], idKey, id, 1)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	var uri bytes.Buffer
	uri.WriteString(path)
	if params != nil && len(params) != 0 {
		uri.WriteString("?")
		uri.WriteString(params.Encode())
	}

	return &CompositeRequest{
		Method:      http.MethodGet,
		URL:         uri.String(),
		ReferenceID: referenceID,
	}
}

//CompositeInsertSObject generates a CompositeRequest in the format of a InsertSObject
func (forceApi *ForceApi) CompositeInsertSObject(in SObject, referenceID string) *CompositeRequest {
	uri := forceApi.apiSObjects[in.ApiName()].URLs[sObjectKey]

	return &CompositeRequest{
		Method:      http.MethodPost,
		URL:         uri,
		Body:        in,
		ReferenceID: referenceID,
	}
}

//CompositeUpdateSObject generates a CompositeRequest in the format of a UpdateSObject
func (forceApi *ForceApi) CompositeUpdateSObject(id string, in SObject, referenceID string) *CompositeRequest {
	uri := strings.Replace(forceApi.apiSObjects[in.ApiName()].URLs[rowTemplateKey], idKey, id, 1)

	return &CompositeRequest{
		Method:      http.MethodPatch,
		URL:         uri,
		Body:        in,
		ReferenceID: referenceID,
	}
}

//CompositeDeleteSObject generates a CompositeRequest in the format of a DeleteSObject
//@Params (obj SObject) only needs an empty object to define SObject type
func (forceApi *ForceApi) CompositeDeleteSObject(id string, obj SObject, referenceID string) *CompositeRequest {
	uri := strings.Replace(forceApi.apiSObjects[obj.ApiName()].URLs[rowTemplateKey], idKey, id, 1)

	return &CompositeRequest{
		Method:      http.MethodDelete,
		URL:         uri,
		ReferenceID: referenceID,
	}
}
