// Package response holds public response, status, buffered transport response,
// and shared response metric payload contracts.
//
// The stable core is:
//   - Status / StatusCoder
//   - Writer
//   - Response / Compressed
//   - Buffered as the transport-ready convenience implementation
//   - Metric payload shapes
//
// Some helper APIs remain public for compatibility, but they should be treated
// as helper surface rather than the governing center of the package.
package response
