package film

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

const (
	HttpURI     = `^/film/(?P<id>[0-9]+)$`
	ParamFilmId = `id`
)

var ResourceDefinitionGet = resource_definition.ResourceDefinition{
	HttpURI:  HttpURI,
	HttpVerb: resource_definition.VerbGet,
	Type:     `FILM_GET`,
	Args: []resource_definition.Arg{
		{
			Key:      ParamFilmId,
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
	HttpVerb: resource_definition.VerbPatch,
	Type:     `FILM_PATCH`,
	Args: []resource_definition.Arg{
		{
			Key:      ParamFilmId,
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
