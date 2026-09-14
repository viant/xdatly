package handler

import "context"

const ReadMetadataKey ValueKey = "readMetadata"

// ReadStep addresses a relation holder and its final row ordinal. Holder is a
// canonical Go holder path, not a SQL alias or a client-supplied field name.
type ReadStep struct {
	Holder string
	Index  int
}

// ReadProjection exposes evidence from one completed typed read. Fields uses
// the final root ordinal followed by optional relation steps. Unknown evidence
// and invalid paths return errors; nil must never mean all fields were loaded.
// RootHolder describes an output envelope; it is empty for direct outputs.
type ReadProjection interface {
	RootHolder() string
	DirectOutput() bool
	Fields(rootOrdinal int, relations ...ReadStep) (FieldSet, error)
}

// ReadOutputProjection optionally addresses additional output-envelope slots
// such as independent derived aggregates. It is separate from ReadProjection
// so existing implementations keep their root-only contract unchanged.
type ReadOutputProjection interface {
	Output(holder string) (ReadProjection, error)
}

// ReadMetadata resolves evidence by the canonical Go field path on the bound
// component input. It describes the successful binding, not a later reread of
// a similarly named view. It contains no database handles or mutable row data.
type ReadMetadata interface {
	Projection(inputField string) (ReadProjection, error)
}

// ReadMetadataConsumer opts a handler/definition into input read evidence.
// Ordinary handlers that do not request it retain their normal read path.
type ReadMetadataConsumer interface {
	RequiresReadMetadata() bool
}

type readMetadataContextKey struct{}
type readMetadataContextValue struct{ value ReadMetadata }

// WithReadMetadata installs invocation-local read-only evidence. A nil value
// shadows inherited evidence so nested components cannot borrow parent fields.
func WithReadMetadata(ctx context.Context, value ReadMetadata) context.Context {
	return context.WithValue(ctx, readMetadataContextKey{}, readMetadataContextValue{value: value})
}

func ReadMetadataFromContext(ctx context.Context) (ReadMetadata, bool) {
	if ctx == nil {
		return nil, false
	}
	value, ok := ctx.Value(readMetadataContextKey{}).(readMetadataContextValue)
	return value.value, ok && value.value != nil
}
