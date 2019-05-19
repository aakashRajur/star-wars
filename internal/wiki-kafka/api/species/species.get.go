package species

import (
	"encoding/json"
	"github.com/aakashRajur/star-wars/pkg/di/kafka-subscription"

	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/species"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourceGet = species.ResourceGet

func GetSpecies(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) kafka_subscription.KafkaSubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		oldPagination := event.Ctx.Value(types.PAGINATION).(types.Pagination)
		result, newPagination, err := species.QuerySelectSpecies(storage, tracker, species.CacheKey, oldPagination)
		if err != nil {
			response.Error = map[string]string{
				`db`: err.Error(),
			}
			err := instance.Emit(response)
			if err != nil {
				logger.Error(err)
			}
			return
		}

		marshaled, err := json.Marshal(*newPagination)
		if err != nil {
			response.Error = map[string]string{
				types.PAGINATION: err.Error(),
			}
			err := instance.Emit(response)
			if err != nil {
				logger.Error(err)
			}
			return
		}

		response.Args = map[string]interface{}{
			types.PAGINATION: string(marshaled),
		}
		response.Data = result
		err = instance.Emit(response)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := kafka.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.Pagination(),
	)

	subscription := kafka.Subscription{
		Topic:   definedTopics[topics.WikiRequestTopic],
		Type:    resourceGet.Type,
		Handler: middlewares(handler),
	}

	return kafka_subscription.KafkaSubscriptionProvider{
		Subscription: &subscription,
	}
}
