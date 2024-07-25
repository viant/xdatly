package validator

import (
	"fmt"
	"sort"
	"strings"
)

type (
	//Violation represent validation violation
	Violation struct {
		Location string      `json:",omitempty"`
		Field    string      `json:",omitempty"`
		Value    interface{} `json:",omitempty"`
		Message  string      `json:",omitempty"`
		Check    string      `json:",omitempty"`
	}

	//Validation represents validation
	Validation struct {
		Violations []*Violation `json:",omitempty"`
		Failed     bool         `json:",omitempty"`
	}

	Violations []*Violation
)

// Sort sorts violations
func (v Violations) Sort() {
	sort.Slice(v, func(i, j int) bool {
		return v[i].Location < v[j].Location
	})
}

// Append appends violation
func (v *Validation) Append(location, field string, value interface{}, check string, msg string) {
	if msg == "" {
		msg = fmt.Sprintf("check '%v' failed on field %v", check, field)
	} else {
		msg = strings.Replace(msg, "$value", fmt.Sprintf("%v", value), 1)
	}
	v.Violations = append(v.Violations, &Violation{
		Location: location,
		Field:    field,
		Message:  msg,
		Check:    check,
		Value:    value,
	})
	v.Failed = len(v.Violations) > 0
}

// AddViolation adds violation
func (v *Validation) AddViolation(location, field string, value interface{}, check string, msg string) {
	v.Append(location, field, value, check, msg)
}

// NewValidation returns a validation
func NewValidation() *Validation {
	return &Validation{Violations: make([]*Violation, 0)}
}
