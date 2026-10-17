package handler

import "context"

// OutputFieldFilter describes the selected public output paths.
type OutputFieldFilter interface{ ExcludePath([]string, string) bool }

// OutputSelection optionally forwards one explicitly delegated child projection.
// The invocation's ComponentInvoker may implement this capability. Callers must
// preserve field paths and publish only the successful child result they return.
type OutputSelection interface {
	CaptureChildOutputSelection(context.Context) context.Context
	SelectedOutputFields(context.Context, any) OutputFieldFilter
	PublishOutputSelection(context.Context, any, OutputFieldFilter)
}
