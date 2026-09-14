package handler

import "context"

// InputSnapshotKey exposes an opt-in capture result through the invocation
// binder. It is runtime-owned and must never be supplied by transport data or
// component providers. Generated adapters own the concrete snapshot type.
const InputSnapshotKey ValueKey = "inputSnapshot"

// InputCapturer optionally captures bound input before Init and InitMCP run.
// CaptureInput must not mutate input. It returns an invocation-local snapshot
// with detached original values/markers, or nil when no capture is needed.
// The runtime treats the result as opaque rather than implementing mutation
// policy or exposing a second binding model.
type InputCapturer[I any] interface {
	CaptureInput(context.Context, *I) (any, error)
}
