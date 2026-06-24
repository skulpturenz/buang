package services

import (
	"reflect"

	dbinterfaces "skulpture/buang/db/interfaces"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
)

const (
	KeyWorkflows string = "workflows"
	KeyPorts     string = "ports"
)

var (
	KeyQueries              = reflect.TypeFor[dbinterfaces.Queries]()
	KeyGorillaSchemaDecoder = reflect.TypeFor[schema.Decoder]()
	KeyGorillaSchemaEncoder = reflect.TypeFor[schema.Encoder]()
	KeyDocker               = reflect.TypeFor[*client.Client]()
)
