package wrappers

import (
	"skulpture/buang/db/interfaces"
	sqlite_models "skulpture/buang/db/sqlite/out"
)

type Webhook sqlite_models.Webhook

var _ interfaces.Webhook = (*Webhook)(nil)

func (w Webhook) GetId() int64 {
	return w.ID
}

func (w Webhook) GetType() int16 {
	return int16(w.Type)
}

func (w Webhook) GetDeleted() bool {
	return w.Deleted != 0
}

func (w Webhook) Unwrap() sqlite_models.Webhook {
	return sqlite_models.Webhook(w)
}

type ProjectWebhook sqlite_models.ProjectWebhook

var _ interfaces.ProjectWebhook = (*ProjectWebhook)(nil)

func (pw ProjectWebhook) GetId() int64 {
	return pw.ID
}

func (pw ProjectWebhook) GetWebhookId() int64 {
	return pw.WebhookID
}

func (pw ProjectWebhook) GetWebhookType() int16 {
	return int16(pw.WebhookType)
}

func (pw ProjectWebhook) GetUrl() string {
	return pw.Url
}

func (pw ProjectWebhook) GetProjectId() int64 {
	return pw.ProjectID
}

func (pw ProjectWebhook) GetDeleted() bool {
	return pw.Deleted != 0
}

func (pw ProjectWebhook) Unwrap() sqlite_models.ProjectWebhook {
	return sqlite_models.ProjectWebhook(pw)
}