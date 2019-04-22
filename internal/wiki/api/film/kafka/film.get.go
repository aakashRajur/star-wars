package kafka

import (
	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourceDefinitionGet = film.ResourceDefinitionGet

func GetFilm(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) di.SubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		args := event.Args
		id, ok := args[film.ParamFilmId].(int)
		if !ok {
			response.Error = map[string]string{
				`id`: `invalid film id`,
			}
			err := instance.Emit(response)
			if err != nil {
				logger.Error(err)
			}
			return
		}

		data, err := film.QuerySelectFilm(storage, tracker, films.CacheKey, id)
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
			resourceDefinitionGet.GetArgValidators(),
			resourceDefinitionGet.GetArgNormalizers(),
			resourceDefinitionGet.DataRequired,
		),
	)

	subscription := kafka.Subscription{
		Topic:   definedTopics[topics.WikiRequestTopic],
		Type:    resourceDefinitionGet.Type,
		Handler: middlewares(handler),
	}

	return di.SubscriptionProvider{
		Subscription: &subscription,
	}
}
