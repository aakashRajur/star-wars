package kafka

import (
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Logger(logger types.Logger) kafka.Middleware {
	return func(next kafka.HandleEvent) kafka.HandleEvent {
		return func(event kafka.Event, kafka *kafka.Kafka) {
			logger.InfoFields(
				map[string]interface{}{
					`TOPIC`:     event.Topic,
					`TYPE`:      event.Type,
					`ID`:        event.Id,
					`SOURCE`:    event.Source,
					`ARGS`:      event.Args,
					`DATA`:      event.Data,
					`TIMESTAMP`: event.Timestamp,
				},
			)

			next(event, kafka)
		}
	}
}
