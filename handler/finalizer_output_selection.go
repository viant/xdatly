package handler

import "context"

// OutputSelectionConsumer opts an output contract into preserving its reader's
// effective projection for its composition finalizer, including typed calls.
// It does not change the projection or public input/output shape.
type OutputSelectionConsumer interface{ RequiresOutputSelection() bool }

// OutputSelectionProvider is an optional capability on the invoker passed to a
// composition finalizer. available distinguishes a successful unrestricted read
// (nil filter, true) from missing or failed reader evidence (nil, false).
// Calls are valid only during the callback. Child reads cannot replace this
// invocation's evidence. A retained immutable filter remains read-only.
// Filter parent paths use encoded JSON names; the leaf field uses its Go name.
type OutputSelectionProvider interface {
	OutputSelection(context.Context) (filter OutputFieldFilter, available bool, err error)
}
