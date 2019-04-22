package kafka

import (
	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/specie"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourceDefinitionPatch = specie.ResourceDefinitionPatch

func PatchSpecie(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) di.SubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		args := event.Args
		id, ok := args[specie.ParamSpecieId].(int)
		if !ok {
			response.Error = map[string]string{
				`id`: `invalid character id`,
			}
			err := instance.Emit(response)
			if err != nil {
				logger.Error(err)
			}
			return
		}

		data := event.Data

		err := specie.QueryUpdateSpecie(storage, tracker, id, data)
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
			resourceDefinitionPatch.GetArgValidators(),
			resourceDefinitionPatch.GetArgNormalizers(),
			true,
		),
		middleware.ValidateData(
			logger,
			definedTopics[topics.WikiResponseTopic],
			specie.SpecieValidation,
			specie.SpecieNormalization,
			resourceDefinitionPatch.DataRequired,
		),
	)

	subscription := kafka.Subscription{
		Topic:   definedTopics[topics.WikiRequestTopic],
		Type:    resourceDefinitionPatch.Type,
		Handler: middlewares(handler),
	}

	return di.SubscriptionProvider{
		Subscription: &subscription,
	}
}
