package module

import (
	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/juju/errors"
	"go.uber.org/fx"
	"math"
)

//noinspection GoSnakeCaseUsage
const (
	KAFKA_SERVICE        = `kafka`
	KAFKA_MAX_PARTITIONS = `KAFKA_MAX_PARTITIONS`
	KAFKA_MAX_REPLICAS   = `KAFKA_MAX_REPLICAS`
)

func GetKafkaConfig(resolver service.Resolver, service service.Service, logger types.Logger, handler types.FatalHandler) kafka.Config {
	brokers, err := resolver.Resolve(KAFKA_SERVICE)
	if err != nil {
		handler.HandleFatal(err)
		return kafka.Config{}
	}
	if len(brokers) < 1 {
		handler.HandleFatal(errors.New(`NO KAFKA BROKERS FOUND`))
		return kafka.Config{}
	}

	maxPartitions := env.GetInt(KAFKA_MAX_PARTITIONS)
	maxReplicas := env.GetInt(KAFKA_MAX_REPLICAS)
	groupId := service.Name
	clientId := service.Hostname

	config := kafka.Config{
		Logger:      logger,
		Brokers:     brokers,
		GroupId:     groupId,
		ClientId:    clientId,
		Replication: int(math.Max(float64(len(brokers)), float64(maxReplicas))),
		Partitions:  int(math.Max(float64(len(brokers)), float64(maxPartitions))),
	}

	return config
}

func GetKafka(config kafka.Config, definedTopics kafka.DefinedTopics, subscriptions []*kafka.Subscription, handler types.FatalHandler) *kafka.Kafka {
	instance, err := kafka.NewInstance(config)
	if err != nil {
		handler.HandleFatal(err)
		return nil
	}

	topics := make([]string, 0)
	for _, topic := range definedTopics {
		topics = append(topics, topic)
	}
	err = instance.CreateTopics(topics...)
	if err != nil {
		handler.HandleFatal(err)
	}

	for _, subscription := range subscriptions {
		err := instance.Subscribe(subscription)
		if err != nil {
			config.Logger.Error(err)
		}
	}

	return instance
}

func GetKafkaProtocol(kafka *kafka.Kafka) types.Protocol {
	return kafka
}

var KafkaModule = fx.Provide(
	GetKafkaConfig,
	GetKafka,
)

var KafkaProtocolModule = fx.Provide(GetKafkaProtocol)
