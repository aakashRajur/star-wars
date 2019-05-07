package planet

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource"
)

const (
	HttpURI       = `/planet/(?P<id>[0-9]+)`
	ParamPlanetId = `id`
)

var ResourceDefinitionGet = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `PLANET_GET`,
	Args: []resource.Arg{
		{
			Key:      ParamPlanetId,
			Type:     resource.TypeInt,
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

var ResourceDefinitionPatch = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbPatch,
	Type:     `PLANET_PATCH`,
	Args: []resource.Arg{
		{
			Key:      ParamPlanetId,
			Type:     resource.TypeInt,
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
