package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

type ResourceDefinitionCompiler struct {
	fx.In
	ResourceDefinitions []resource_definition.ResourceDefinition `group:"resource_definitions"`
}

type ResourceDefinitionProvider struct {
	fx.Out
	ResourceDefinition resource_definition.ResourceDefinition `group:"resource_definitions"`
}
