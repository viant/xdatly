// Package handler holds the stabilized public custom-handler contract surface
// for Datly 1.0.
//
// Its core is intentionally small:
//   - Func[I,O] for simple pure handlers
//   - Contract[I,O] for explicit runtime-aware handlers
//   - Binder plus InputKey for invocation-scoped value access
//   - Session as the reduced binder/response facade; logger, validator, message
//     bus, and data services are binder-scoped capabilities rather than session APIs
//   - generic session and pre-invoke hook context carriers
//   - plain public lifecycle hooks
//
// It must not regrow Datly runtime internals such as SQL execution,
// transaction/state managers, route execution, or transport-specific hooks.
package handler
