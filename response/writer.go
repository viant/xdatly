package response

// Writer is the minimal handler-visible response/status facade.
// Ownership of final response/meta/status shaping stays with runtime.
type Writer interface {
	StatusCode() int
	SetStatusCode(code int)
	AddError(err error)
	AddMetric(metric *Metric)
	Metrics() Metrics
}
