package handler

import (
	"context"
	"github.com/viant/bindly/locator"
)

// ComponentInvokerKey identifies the invocation-scoped composition capability.
const ComponentInvokerKey ValueKey = "component_invoker"

// ComponentTarget selects an exact component and route. This is trusted
// application configuration, never a target selected by transport input.
type ComponentTarget struct {
	Scope  string
	Name   string
	Method string
	Path   string
}

// ComponentRequest composes a child through the runtime's canonical binding
// and execution path. Providers overlay the inherited request scope. Calls
// retain the caller's buffered writes and transaction completion owner.
type ComponentRequest struct {
	Target     ComponentTarget
	Input      any
	ExtraInput *InputSeed         `json:"-"`
	Providers  []locator.Provider `json:"-"`
	Completion func(Outcome)      `json:"-"`
}

// ComponentInvoker is implemented by the runtime and injected into application
// handlers. Invoke must be awaited; it does not independently flush or commit
// a caller-owned unit of work.
type ComponentInvoker interface {
	Invoke(context.Context, ComponentRequest) (any, error)
}
