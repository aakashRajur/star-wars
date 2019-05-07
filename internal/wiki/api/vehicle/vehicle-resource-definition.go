package vehicle

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource"
)

const (
	HttpURI        = `/vehicle/(?P<id>[0-9]+)`
	ParamVehicleId = `id`
)

var ResourceDefinitionGet = resource.Definition{
	HttpURI:  HttpURI,
	HttpVerb: http.VerbGet,
	Type:     `VEHICLE_GET`,
	Args: []resource.Arg{
		{
			Key:      ParamVehicleId,
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
	Type:     `VEHICLE_PATCH`,
	Args: []resource.Arg{
		{
			Key:      ParamVehicleId,
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
