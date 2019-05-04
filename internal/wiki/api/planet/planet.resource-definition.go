package planet

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

const (
	HttpURI       = `/planet/(?P<id>[0-9]+)`
	ParamPlanetId = `id`
)

var ResourceDefinitionGet = resource_definition.ResourceDefinition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `PLANET_GET`,
	Args: []resource_definition.Arg{
		{
			Key:      ParamPlanetId,
			Type:     resource_definition.TypeInt,
			Required: true,
		},
	},
	DataRequired: false,
}

func ProvideResourceDefinitionGet() di.ResourceDefinitionProvider {
	return di.ResourceDefinitionProvider{
		ResourceDefinition: ResourceDefinitionGet,
	}
}

var ResourceDefinitionPatch = resource_definition.ResourceDefinition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbPatch,
	Type:     `PLANET_PATCH`,
	Args: []resource_definition.Arg{
		{
			Key:      ParamPlanetId,
			Type:     resource_definition.TypeInt,
			Required: true,
		},
	},
	DataRequired: true,
}

func ProvideResourceDefinitionPatch() di.ResourceDefinitionProvider {
	return di.ResourceDefinitionProvider{
		ResourceDefinition: ResourceDefinitionPatch,
	}
}
