package workersshared

import (
	"skulpture/buang/app"
	dbinterfaces "skulpture/buang/db/interfaces"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
)

type WorkflowServices struct {
	Queries              *dbinterfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               *client.Client
}

func (ws WorkflowServices) ToAppServices() app.ApplicationServices {
	return app.ApplicationServices{
		Queries:              ws.Queries,
		GorillaSchemaDecoder: ws.GorillaSchemaDecoder,
		GorillaSchemaEncoder: ws.GorillaSchemaEncoder,
		Docker:               ws.Docker,
	}
}
