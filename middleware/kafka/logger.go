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
					`topic`:     event.Topic,
					`type`:      event.Type,
					`id`:        event.Id,
					`source`:    event.Source,
					`args`:      event.Args,
					`data`:      event.Data,
					`timestamp`: event.Timestamp,
				},
			)

			next(event, kafka)
		}
	}
}
