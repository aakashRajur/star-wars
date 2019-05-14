package film

import (
	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourceGet = film.ResourceGet

func GetFilm(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) di.KafkaSubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		args := event.Args
		id := args[film.ParamFilmId].(int)

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
			film.ArgValidation,
			film.ArgNormalization,
			true,
		),
	)

	subscription := kafka.Subscription{
		Topic:   definedTopics[topics.WikiRequestTopic],
		Type:    resourceGet.Type,
		Handler: middlewares(handler),
	}

	return di.KafkaSubscriptionProvider{
		Subscription: &subscription,
	}
}
