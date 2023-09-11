package async

const (
	InvocationTypeEvent     InvocationType = "event"
	InvocationTypeUndefined InvocationType = ""
)

type (
	InvocationType    string
	invocationTypeKey string
)

// InvocationTypeKey defines context key
var InvocationTypeKey = "invocation"
