package film

import (
	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/di/kafka-subscription"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourcePatch = film.ResourcePatch

func PatchFilm(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) kafka_subscription.KafkaSubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		args := event.Args
		id := args[film.ParamFilmId].(int)

		data := event.Data.(map[string]interface{})

		err := film.QueryUpdateFilm(storage, tracker, id, data)

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
			film.ArgValidation,
			film.ArgNormalization,
			true,
		),
		middleware.ValidateData(
			logger,
			definedTopics[topics.WikiResponseTopic],
			film.BodyValidation,
			film.BodyNormalization,
			true,
		),
	)

	subscription := kafka.Subscription{
		Topic:   definedTopics[topics.WikiRequestTopic],
		Type:    resourcePatch.Type,
		Handler: middlewares(handler),
	}

	return kafka_subscription.KafkaSubscriptionProvider{
		Subscription: &subscription,
	}
}
