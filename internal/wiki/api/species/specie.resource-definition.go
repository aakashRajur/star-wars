package species

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	HttpURI = `^/species$`
)

var ResourceDefinitionGet = resource_definition.ResourceDefinition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `SPECIES_GET`,
	Args: []resource_definition.Arg{
		{
			Key:      types.QueryPaginationId,
			Type:     resource_definition.TypeInt,
			Required: false,
		},
		{
			Key:      types.QueryLimit,
			Type:     resource_definition.TypeInt,
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
