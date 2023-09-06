package async

const (
	NotificationMethodStorage    NotificationMethod = "Storage"
	NotificationMethodMessageBus NotificationMethod = "MessageBus"
	NotificationMethodChannel    NotificationMethod = "Channel"
	NotificationMethodUndefined  NotificationMethod = ""
)

type (
	NotificationMethod string

	Notification struct {
		Method      NotificationMethod
		Destination string
	}
)
