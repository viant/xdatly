package handler

// Conflict reports a failed comparison with the authorized Previous projection.
// It is a validation result; it does not imply database locking or an atomic
// compare-and-write operation.
type Conflict struct {
	Entity string
	Field  string
	Reason string
}

func (e *Conflict) Error() string {
	return "mutation conflict: " + e.Entity + "." + e.Field + ": " + e.Reason
}

// StatusCode exposes the validation conflict through the standard response contract.
func (e *Conflict) StatusCode() int { return 409 }
