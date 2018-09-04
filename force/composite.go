package force

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
	Body           interface{}       `force:"body,omitempty"`
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

	return
}

//Add adds a request object to the CompositeRequests struct so that when it is executed, all containing requests are processed
func (requests *CompositeRequests) Add(request *CompositeRequest) {
	requests.CompositeRequest = append(requests.CompositeRequest, *request)
	return
}
