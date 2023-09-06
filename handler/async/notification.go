package async

const (
	NotificationMethodS3        NotificationMethod = "S3"
	NotificationMethodSQS       NotificationMethod = "SQS"
	NotificationMethodUndefined NotificationMethod = ""
)

type (
	NotificationMethod string

	Notification struct {
		HandlerType NotificationMethod
		Destination string
	}
)
