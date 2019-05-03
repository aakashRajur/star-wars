package types

import (
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

type ResourceRegistration interface {
	Register(definition resource_definition.ResourceDefinition) error
	Unregister(definition resource_definition.ResourceDefinition) error
}
