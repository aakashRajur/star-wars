package vehicle

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

const (
	HttpURI        = `^/vehicle/(?P<id>[0-9]+)$`
	ParamVehicleId = `id`
)

var ResourceDefinitionGet = resource_definition.ResourceDefinition{
	HttpURI:  HttpURI,
	HttpVerb: resource_definition.VerbGet,
	Type:     `VEHICLE_GET`,
	Args: []resource_definition.Arg{
		{
			Key:      ParamVehicleId,
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
	Type:     `VEHICLE_PATCH`,
	Args: []resource_definition.Arg{
		{
			Key:      ParamVehicleId,
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
