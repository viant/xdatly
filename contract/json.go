package contract

import "github.com/viant/x/shape"

// JSONEncoder is an optional capability of a configured encoder. It receives
// immutable field metadata and the original linked value, never a replacement
// receiver type. Implementations retain their existing JSON ownership rules.
type JSONEncoder interface {
	MarshalWithRepresentation(value any, metadata *shape.Representation) ([]byte, error)
}

// JSONDecoder is the corresponding optional capability of a configured decoder.
// The destination retains its original linked type and presence-marker layout.
type JSONDecoder interface {
	UnmarshalWithRepresentation(data []byte, destination any, metadata *shape.Representation) error
}
