package film

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource"
)

const (
	HttpURI     = `/film/(?P<id>[0-9]+)`
	ParamFilmId = `id`
)

var ResourceDefinitionGet = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `FILM_GET`,
	Args: []resource.Arg{
		{
			Key:      ParamFilmId,
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
	Type:     `FILM_PATCH`,
	Args: []resource.Arg{
		{
			Key:      ParamFilmId,
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
