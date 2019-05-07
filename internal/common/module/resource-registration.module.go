package module

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/resource"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetResourceRegistrationModule(
	logger types.Logger,
	handler types.FatalHandler,
	lifecycle fx.Lifecycle,
	endpoint types.Endpoint,
	protocol types.Protocol,
	registrar resource.Registree,
	definitions []resource.Definition,
) {
	accessUriPrefix := env.GetString(`ACCESS_URI_PREFIX`)
	if accessUriPrefix == `` {
		handler.HandleFatal(errors.New(`ACCESS_URI_PREFIX WAS NOT SET`))
		return
	}

	compiled := make([]resource.Definition, 0)
	for _, definition := range definitions {
		enriched := definition.Copy()
		enriched.AccessURI = fmt.Sprintf(`%s%s`, accessUriPrefix, enriched.HttpURI)
		enriched.Source = string(endpoint)
		enriched.Protocol = protocol.Name()
		enriched.Endpoint = string(endpoint)
		compiled = append(compiled, enriched)
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				for _, enriched := range compiled {
					err := registrar.Register(enriched)
					if err != nil {
						logger.Fatal(fmt.Sprintf(`UNABLE TO REGISTER RESOURCE: %s`, enriched.Type))
						continue
					}
					logger.Info(fmt.Sprintf(`REGSITERED RESOURCE: %s`, enriched.Type))
				}
				return nil
			},
			OnStop: func(context.Context) error {
				for _, enriched := range compiled {
					err := registrar.Unregister(enriched)
					if err != nil {
						logger.Fatal(fmt.Sprintf(`UNABLE TO UNREGISTER RESOURCE: %s`, enriched.Type))
						continue
					}
					logger.Info(fmt.Sprintf(`UNREGSITERED RESOURCE: %s`, enriched.Type))
				}
				return nil
			},
		},
	)
}

var ResourceRegistrationModule = fx.Invoke(GetResourceRegistrationModule)
