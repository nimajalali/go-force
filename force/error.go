package force

import (
	"fmt"
)

// Custom Error to handle salesforce api responses.
type ApiErrors []*ApiError

type ApiError struct {
	Fields           []string `json:"fields,omitempty" force:"fields,omitempty"`
	Message          string   `json:"message,omitempty" force:"message,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty" force:"errorCode,omitempty"`
	ErrorName        string   `json:"error,omitempty" force:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty" force:"error_description,omitempty"`
}

func (e ApiErrors) Error() string {
	return fmt.Sprintf("%#v", e.Errors)
}

func (e ApiErrors) Errors() []string {
	eArr := make([]string, len(e))
	for i, err := range e {
		eArr[i] = err.Error()
	}
	return eArr
}

func (e ApiErrors) Validate() bool {
	if len(e) != 0 {
		return true
	}

	return false
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%#v", e)
}

func (e ApiError) Validate() bool {
	if len(e.Fields) != 0 || len(e.Message) != 0 || len(e.ErrorCode) != 0 || len(e.ErrorName) != 0 || len(e.ErrorDescription) != 0 {
		return true
	}

	return false
}
