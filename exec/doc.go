// Package exec holds the reduced public execution-context contract for Datly
// runtimes.
//
// The stable core is the generic execution/observability payload:
//   - Context
//   - ContextKey
//   - WithContext / GetContext
//   - New(...) plus its generic options
//
// Explicit HTTP construction and logging-preparation helpers may remain public
// for compatibility, but they are not the governing center of the package.
package exec
