package xdatly

// Component is the minimal public authoring primitive for Datly 1.0.
// It is intentionally small and does not become the canonical IR by itself.
type Component[I any, O any] struct {
	Input  I
	Output O
}
