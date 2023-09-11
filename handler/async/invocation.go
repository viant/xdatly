package async

const (
	InvocationTypeEvent     = "event"
	InvocationTypeUndefined = ""
)

type (
	InvocationType    string
	invocationTypeKey string
)

// InvocationTypeKey defines context key
var InvocationTypeKey = "invocation"
