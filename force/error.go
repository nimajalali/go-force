package force

import (
	"fmt"
	"strings"
)

// APIErrors is an errors wrapper for Salesforce API responses
type APIErrors []*APIError

// APIError represents a Salesforce API error response
type APIError struct {
	Fields           []string `json:"fields,omitempty" force:"fields,omitempty"`
	Message          string   `json:"message,omitempty" force:"message,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty" force:"errorCode,omitempty"`
	ErrorName        string   `json:"error,omitempty" force:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty" force:"error_description,omitempty"`
}

func (e APIErrors) Error() string {
	return e.String()
}

func (e APIErrors) String() string {
	s := make([]string, len(e))
	for i, err := range e {
		s[i] = err.String()
	}

	return strings.Join(s, "\n")
}

func (e APIErrors) Validate() bool {
	for _, err := range e {
		if err.ErrorCode != "" {
			return true
		}
	}

	return false
}

func (e APIError) Error() string {
	return e.String()
}

func (e APIError) String() string {
	return fmt.Sprintf("%#v", e)
}

func (e APIError) Validate() bool {
	if len(e.Fields) != 0 || len(e.Message) != 0 || len(e.ErrorCode) != 0 || len(e.ErrorName) != 0 || len(e.ErrorDescription) != 0 {
		return true
	}

	return false
}
