package handler

// ValidationIssue is read-only schema evidence for aggregate business validation.
// Raw checked values are deliberately excluded: they may reference live input.
// The native Validation retains that evidence for the normal error finalizer.
type ValidationIssue struct {
	Location string
	Field    string
	Message  string
	Check    string
}

// ValidationReport collects one invocation's initial schema and business
// violations. Ordinary schema violations do not suppress aggregate validation.
// Operational validator errors do; they are not converted into report entries.
// Implementations return detached schema evidence and retain all native failures.
// The callback must not retain the report or use it concurrently.
type ValidationReport interface {
	SchemaFailed() bool
	SchemaViolations() []ValidationIssue
	Add(Violation)
}
