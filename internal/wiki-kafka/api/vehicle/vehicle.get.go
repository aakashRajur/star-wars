package vehicle

import (
	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles"
	middleware "github.com/aakashRajur/star-wars/middleware/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

var resourceGet = vehicle.ResourceGet

func GetVehicle(storage types.Storage, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) di.KafkaSubscriptionProvider {
	handler := func(event kafka.Event, instance *kafka.Kafka) {
		response := kafka.Event{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  event.Type,
			Id:    event.Id,
		}

		args := event.Args
		id := args[vehicle.ParamVehicleId].(int)

		data, err := vehicle.QuerySelectVehicle(storage, tracker, vehicles.CacheKey, id)
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
			vehicle.ArgValidation,
			vehicle.ArgNormalization,
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
