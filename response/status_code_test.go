package response

import (
	"errors"
	"testing"
)

type testStatusError struct {
	code int
}

func (e testStatusError) Error() string   { return "boom" }
func (e testStatusError) StatusCode() int { return e.code }

func TestErrorStatusCode(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		fallback int
		expect   int
	}{
		{name: "nil", fallback: 500, expect: 500},
		{name: "direct", err: testStatusError{code: 409}, fallback: 500, expect: 409},
		{name: "wrapped", err: errors.Join(errors.New("wrapped"), testStatusError{code: 422}), fallback: 500, expect: 422},
		{name: "zero code", err: testStatusError{code: 0}, fallback: 500, expect: 500},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := ErrorStatusCode(testCase.err, testCase.fallback); actual != testCase.expect {
				t.Fatalf("expected %d, got %d", testCase.expect, actual)
			}
		})
	}
}
