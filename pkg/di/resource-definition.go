package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/resource"
)

type ResourceDefinitionCompiler struct {
	fx.In
	ResourceDefinitions []resource.Definition `group:"resource_definitions"`
}

type ResourceDefinitionProvider struct {
	fx.Out
	ResourceDefinition resource.Definition `group:"resource_definitions"`
}
