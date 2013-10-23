package force

import (
	"fmt"
)

const (
	resourcesUri = "/services/data/%v"
)

var ApiResources map[string]string

func getApiResources() error {
	uri := fmt.Sprintf(resourcesUri, apiVersion)

	ApiResources = make(map[string]string)
	return get(uri, nil, &ApiResources)
}

// Custom Error to handle salesforce api responses.
type ApiErrors []ApiError

type ApiError struct {
	Fields           []string `json:"fields,omitempty" force:"fields,omitempty"`
	Message          string   `json:"message,omitempty" force:"message,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty" force:"errorCode,omitempty"`
	ErrorName        string   `json:"error,omitempty" force:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty" force:"error_description,omitempty"`
}

func (e ApiErrors) Error() string {
	return fmt.Sprintf("%#v", e)
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%#v", e)
}
