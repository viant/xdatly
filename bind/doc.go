// Package bind holds public binder-scoped provider, scope, and typed lookup
// contracts.
//
// It exists for stable DI/value-resolution contracts only. It intentionally
// stops short of promoting Datly's current reflection/type-specific injection
// logic, because that logic is still runtime-colored rather than a proven
// cross-runtime public contract.
package bind
