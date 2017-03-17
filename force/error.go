package force

import (
	"fmt"
	"strings"
)

// APIErrors is a list of API errors used to handle salesforce API responses.
type APIErrors []*APIError

// APIError reprensents a SalesForce API error.
type APIError struct {
	Fields           []string `json:"fields,omitempty" force:"fields,omitempty"`
	Message          string   `json:"message,omitempty" force:"message,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty" force:"errorCode,omitempty"`
	ErrorName        string   `json:"error,omitempty" force:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty" force:"error_description,omitempty"`
}

// Error returns the string representation for an APIErrors.
func (e APIErrors) Error() string {
	return e.String()
}

// String formats the fields in an APIErrors.
func (e APIErrors) String() string {
	s := make([]string, len(e))
	for i, err := range e {
		s[i] = err.String()
	}

	return strings.Join(s, "\n")
}

// Validate validates an APIErrors.
func (e APIErrors) Validate() bool {
	return len(e) != 0
}

// Error returns the string representation for an APIError.
func (e APIError) Error() string {
	return e.String()
}

// String formats the fields in an APIError.
func (e APIError) String() string {
	return fmt.Sprintf("%#v", e)
}

// Validate validates an APIError.
func (e APIError) Validate() bool {
	if len(e.Fields) != 0 || len(e.Message) != 0 || len(e.ErrorCode) != 0 || len(e.ErrorName) != 0 || len(e.ErrorDescription) != 0 {
		return true
	}

	return false
}
