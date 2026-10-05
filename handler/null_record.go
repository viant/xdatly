package handler

import "fmt"

// RootNullRecordError is an operational structural error from initial native
// validation under the explicit root-null policy. It is not a business violation.
// Location identifies the original canonical collection position; Cause retains
// private diagnostics. Applications may project their existing error envelope.
type RootNullRecordError struct {
	Location string
	Cause    error
}

func (e *RootNullRecordError) Error() string {
	return fmt.Sprintf("null root record at %s: %v", e.Location, e.Cause)
}
func (e *RootNullRecordError) Unwrap() error { return e.Cause }
