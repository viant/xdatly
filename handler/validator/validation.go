package validator

import (
	"fmt"
	"strings"
)

type (
	//Violation represent validation violation
	Violation struct {
		Location string
		Field    string
		Value    interface{}
		Message  string
		Check    string
	}

	//Validation represents validation
	Validation struct {
		Violations []*Violation
		Failed     bool
	}
)

//Append appends violation
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

//AddViolation adds violation
func (v *Validation) AddViolation(location, field string, value interface{}, check string, msg string) {
	v.Append(location, field, value, check, msg)
}
