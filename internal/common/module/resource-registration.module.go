package module

import (
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
	"github.com/aakashRajur/star-wars/pkg/types"
	"go.uber.org/fx"
)

func GetResourceRegistrationModule(instanceId types.InstanceId, protocol types.Protocol, definitions []resource_definition.ResourceDefinition) {

}

var ResourceRegistrationModule = fx.Invoke(GetResourceRegistrationModule)
