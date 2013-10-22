package force

import (
	"fmt"
)

// Custom Error to handle salesforce api responses.
type ApiError struct {
	Fields           []string `json:"fields,omitempty"`
	Message          string   `json:"message,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty"`
	ErrorName        string   `json:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("force error: %#v", e)
}
