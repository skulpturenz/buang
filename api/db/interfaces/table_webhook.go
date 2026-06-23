package interfaces

type Webhook interface {
	GetId() int64
	GetType() int16
	GetDeleted() bool
}

type ProjectWebhook interface {
	GetId() int64
	GetWebhookId() int64
	GetWebhookType() int16
	GetUrl() string
	GetProjectId() int64
	GetDeleted() bool
}