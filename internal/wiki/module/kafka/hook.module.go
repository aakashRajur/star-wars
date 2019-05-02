package kafka

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

func GetKafkaHook(config kafka.Config, definedTopics kafka.DefinedTopics, definitions []resource_definition.ResourceDefinition) kafka.Hook {
	hook := kafka.Hook{
		OnStart: func(instance *kafka.Kafka) {
			/*logger := config.Logger
			config := instance.Config

			for _, each := range definitions {
				event := kafka.Event{
					Topic:  definedTopics[topics.ResourceDiscoveryTopic],
					Id:     each.Type,
					Type:   resource_definition.TypeResourceDefinition,
					Source: config.ClientId,
					Data:   each.GetMap(),
				}

				err := instance.Emit(event)
				if err != nil {
					logger.Error(err)
				}
			}*/
		},
	}

	return hook
}

var HookModule = fx.Provide(GetKafkaHook)
