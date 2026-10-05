package handler

import "fmt"

// NestedNullRecordError retains the original location and private cause of a
// null element in an explicitly opted-in writable collection relation.
// It is an operational structural error, not a validation violation.
type NestedNullRecordError struct {
	Location string
	Cause    error
}

func (e *NestedNullRecordError) Error() string {
	return fmt.Sprintf("null nested record at %s: %v", e.Location, e.Cause)
}
func (e *NestedNullRecordError) Unwrap() error { return e.Cause }
