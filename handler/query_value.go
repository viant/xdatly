package handler

// TypedQueryValue carries an already decoded query value for a trusted internal
// invocation. HTTP query occurrences continue to use ordinary wire decoding.
type TypedQueryValue struct{ Value any }
