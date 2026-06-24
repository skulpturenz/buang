package wrappers

import (
	"skulpture/buang/db/interfaces"
	pg_models "skulpture/buang/db/pg/out"
)

type Webhook pg_models.Webhook

var _ interfaces.Webhook = (*Webhook)(nil)

func (w Webhook) GetId() int64 {
	return w.ID
}

func (w Webhook) GetType() int16 {
	return w.Type
}

func (w Webhook) GetDeleted() bool {
	return w.Deleted
}

func (w Webhook) Unwrap() pg_models.Webhook {
	return pg_models.Webhook(w)
}

type ProjectWebhook pg_models.ProjectWebhook

var _ interfaces.ProjectWebhook = (*ProjectWebhook)(nil)

func (pw ProjectWebhook) GetId() int64 {
	return pw.ID
}

func (pw ProjectWebhook) GetWebhookId() int64 {
	return pw.WebhookID
}

func (pw ProjectWebhook) GetWebhookType() int16 {
	return pw.WebhookType
}

func (pw ProjectWebhook) GetUrl() string {
	return pw.Url
}

func (pw ProjectWebhook) GetProjectId() int64 {
	return pw.ProjectID
}

func (pw ProjectWebhook) GetDeleted() bool {
	return pw.Deleted
}

func (pw ProjectWebhook) Unwrap() pg_models.ProjectWebhook {
	return pg_models.ProjectWebhook(pw)
}