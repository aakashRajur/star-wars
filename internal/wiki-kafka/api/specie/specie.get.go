package specie

import (
	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/specie"
	"github.com/aakashRajur/star-wars/internal/wiki/api/species"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/di/kafka-subscription"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourceGet = specie.ResourceGet

func GetSpecie(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) kafka_subscription.KafkaSubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		args := event.Args
		id := args[specie.ParamSpecieId].(int)

		data, err := specie.QuerySelectSpecie(storage, tracker, species.CacheKey, id)
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

		response.Data = data
		err = instance.Emit(response)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := kafka.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.ValidateArgs(
			logger,
			definedTopics[topics.WikiResponseTopic],
			specie.ArgValidation,
			specie.ArgNormalization,
			true,
		),
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
