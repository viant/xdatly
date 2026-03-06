package xdatly

type Component[I any, O any] struct {
	Inout  I
	Output O
}
