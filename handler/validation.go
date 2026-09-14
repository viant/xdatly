package handler

import (
	"net/http"
	"strings"
)

// Violation describes a failed validation rule without owning validation execution.
type Violation struct {
	Location string `json:"location,omitempty"`
	Field    string `json:"field,omitempty"`
	Message  string `json:"message,omitempty"`
	Check    string `json:"check,omitempty"`
}

// Validation is the typed result shared by Go and template validation consumers.
// Either Failed or a nonempty Violations list represents failure.
type Validation struct {
	Failed     bool         `json:"failed,omitempty"`
	Violations []*Violation `json:"violations,omitempty"`
	Code       int          `json:"-"`
}

// StatusCode defaults failed validation to 422; applications may select an
// explicit code, such as 401 for an authorization-related validation result.
func (v *Validation) StatusCode() int {
	if v == nil || (!v.Failed && len(v.Violations) == 0) {
		return 0
	}
	if v.Code != 0 {
		return v.Code
	}
	return http.StatusUnprocessableEntity
}
func (v *Validation) ResponseBody() any { return v }

// Err preserves the typed result for errors.As and returns nil for success.
func (v *Validation) Err() error {
	if v == nil || !v.Failed && len(v.Violations) == 0 {
		return nil
	}
	return v
}

func (v *Validation) Error() string { return strings.Join(v.Messages(), "; ") }

// Messages returns detached human-readable diagnostics. Raw input values are
// not included automatically in errors or protocol responses.
func (v *Validation) Messages() []string {
	if v == nil {
		return nil
	}
	var result []string
	for _, violation := range v.Violations {
		message := "validation failed"
		if violation != nil && violation.Message != "" {
			message = violation.Message
		}
		result = append(result, message)
	}
	if v.Failed && len(result) == 0 {
		result = append(result, "validation failed")
	}
	return result
}
