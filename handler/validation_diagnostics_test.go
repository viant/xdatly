package handler

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCheckedValueIsPrivateByDefault(t *testing.T) {
	for _, testCase := range []struct {
		description string
		input       any
		expected    bool
	}{
		{description: "checked zero", input: 0, expected: true},
		{description: "checked false", input: false, expected: true},
		{description: "checked null", input: nil, expected: true},
		{description: "checked sensitive value", input: "private-value", expected: true},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			value := &Violation{Location: "Rows[1].Name", Field: "Name", Check: "required", CheckedValue: testCase.input, HasCheckedValue: testCase.expected}
			encoded, err := json.Marshal(&Validation{Failed: true, Violations: []*Violation{value}})
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(encoded), "private-value") || strings.Contains(string(encoded), "CheckedValue") || strings.Contains(string(encoded), "checkedValue") {
				t.Fatalf("private checked evidence exposed: %s", encoded)
			}
			if value.HasCheckedValue != testCase.expected {
				t.Fatal("checked null evidence lost")
			}
		})
	}
}
