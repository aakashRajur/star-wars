package species

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource"
	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	HttpURI = `/species`
)

var ResourceDefinitionGet = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `SPECIES_GET`,
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
