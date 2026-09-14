package handler

import (
	"errors"
	"testing"
)

func TestValidationOutcome(t *testing.T) {
	for _, tt := range []struct {
		name   string
		result *Validation
		failed bool
	}{
		{"nil", nil, false}, {"empty", &Validation{}, false},
		{"flag", &Validation{Failed: true}, true},
		{"violations without flag", &Validation{Violations: []*Violation{{Field: "Name", Message: "name required"}}}, true},
		{"nil violation", &Validation{Violations: []*Violation{nil}}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.Err()
			if (err != nil) != tt.failed {
				t.Fatalf("error=%v", err)
			}
			if tt.failed {
				var actual *Validation
				if !errors.As(err, &actual) || actual != tt.result || err.Error() == "" {
					t.Fatalf("typed failure lost: %v", err)
				}
			}
		})
	}
}
