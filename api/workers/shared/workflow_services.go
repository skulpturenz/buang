package workersshared

import (
	"skulpture/buang/app"
	dbinterfaces "skulpture/buang/db/interfaces"
	"skulpture/buang/ports"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
)

type WorkflowServices struct {
	Queries              *dbinterfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               *client.Client
	Ports                ports.Ports
}

func (ws WorkflowServices) ToAppServices() app.ApplicationServices {
	return app.ApplicationServices{
		Queries:              ws.Queries,
		GorillaSchemaDecoder: ws.GorillaSchemaDecoder,
		GorillaSchemaEncoder: ws.GorillaSchemaEncoder,
		Docker:               ws.Docker,
		Ports:                ws.Ports,
	}
}
