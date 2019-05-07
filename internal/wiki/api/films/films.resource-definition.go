package films

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource"
	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	HttpURI = `/films`
)

var ResourceDefinitionGet = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `FILMS_GET`,
	Args: []resource.Arg{
		{
			Key:      types.QueryPaginationId,
			Type:     resource.TypeInt,
			Required: false,
		},
		{
			Key:      types.QueryLimit,
			Type:     resource.TypeInt,
			Required: false,
		},
	},
	DataRequired: false,
}

func ProvideResourceDefinitionGet() di.ResourceDefinitionProvider {
	return di.ResourceDefinitionProvider{
		ResourceDefinition: ResourceDefinitionGet,
	}
}
