package differ

import "context"

// Differ compares values without mutating them. Comparison failures are explicit;
// native comparison errors are returned and retained in ChangeLog.
type Differ interface {
	Diff(ctx context.Context, from, to any, options ...Option) (*ChangeLog, error)
}
