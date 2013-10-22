package force

import (
	"fmt"
)

// Custom Error to handle salesforce api responses.
type ApiError struct {
	Fields    []string `json:"fields,omitempty"`
	Message   string   `json:"message,omitempty"`
	ErrorCode string   `json:"errorCode,omitempty"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("force error: %#v", e)
}
