// Package tracing holds public tracing payload contracts plus a small set of
// compatibility helpers.
//
// The core public surface is the payload model:
//   - Trace
//   - Span
//   - SpanStatus
//   - ResourceInfo
//
// Transport-colored helpers such as HTTP status mapping may remain public for
// compatibility, but they are not the governing center of the package.
package tracing
