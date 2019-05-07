package character

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource"
)

const (
	HttpURI          = `/character/(?P<id>[0-9]+)`
	ParamCharacterId = `id`
)

var ResourceDefinitionGet = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `CHARACTER_GET`,
	Args: []resource.Arg{
		{
			Key:      ParamCharacterId,
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
	Type:     `CHARACTER_PATCH`,
	Args: []resource.Arg{
		{
			Key:      ParamCharacterId,
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
