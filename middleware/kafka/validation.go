package kafka

import (
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate"
)

func ValidateArgs(logger types.Logger, responseTopic string, validators map[string][]types.Validator, normalizors map[string]types.Normalizor, strict bool) kafka.Middleware {
	return func(next kafka.HandleEvent) kafka.HandleEvent {
		return func(event kafka.Event, instance *kafka.Kafka) {
			response := kafka.Event{
				Topic: responseTopic,
				Type:  event.Type,
				Id:    event.Id,
			}
			args := event.Args
			if args == nil {
				if strict {
					response.Error = map[string]string{
						`args`: `args is required`,
					}
					err := instance.Emit(response)
					if err != nil {
						logger.Error(err)
					}
					return
				}
			}

			err := validate.ValidateMapped(validators, args)
			if err != nil {
				response.Error = err
				err := instance.Emit(response)
				if err != nil {
					logger.Error(err)
				}
				return
			}

			if normalizors == nil {
				next(event, instance)
				return
			}

			normalized, err := validate.NormalizeMapped(normalizors, args)
			if err != nil {
				response.Error = err
				err := instance.Emit(response)
				if err != nil {
					logger.Error(err)
				}
				return
			}

			event.Args = normalized
			next(event, instance)
		}
	}
}

func ValidateData(logger types.Logger, eventType string, validators map[string][]types.Validator, normalizors map[string]types.Normalizor, dataRequired bool) kafka.Middleware {
	return func(next kafka.HandleEvent) kafka.HandleEvent {
		return func(event kafka.Event, instance *kafka.Kafka) {
			response := kafka.Event{
				Topic: event.Topic,
				Type:  eventType,
				Id:    event.Id,
			}
			data := event.Data
			if data == nil {
				if dataRequired {
					response.Error = map[string]string{
						`data`: `data is required`,
					}
					err := instance.Emit(response)
					if err != nil {
						logger.Error(err)
					}
					return
				}
			}

			err := validate.ValidateMapped(validators, data)
			if err != nil {
				response.Error = err
				err := instance.Emit(response)
				if err != nil {
					logger.Error(err)
				}
				return
			}

			if normalizors == nil {
				next(event, instance)
				return
			}

			normalized, err := validate.NormalizeMapped(normalizors, data)
			if err != nil {
				response.Error = err
				err := instance.Emit(response)
				if err != nil {
					logger.Error(err)
				}
				return
			}

			event.Data = normalized
			next(event, instance)
		}
	}
}
