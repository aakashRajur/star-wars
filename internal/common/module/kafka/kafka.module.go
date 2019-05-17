package kafka

import (
	"math"

	"github.com/juju/errors"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

//noinspection GoSnakeCaseUsage
const (
	KAFKA_SERVICE        = `kafka`
	KAFKA_MAX_PARTITIONS = `KAFKA_MAX_PARTITIONS`
	KAFKA_MAX_REPLICAS   = `KAFKA_MAX_REPLICAS`
)

func GetKafkaBrokers(resolver service.Resolver, handler types.FatalHandler) kafka.Brokers {
	brokers, err := resolver.Resolve(KAFKA_SERVICE)
	if err != nil {
		handler.HandleFatal(err)
		return []string{}
	}
	if len(brokers) < 1 {
		handler.HandleFatal(errors.New(`NO KAFKA BROKERS FOUND`))
		return []string{}
	}
	return brokers
}

func GetKafkaConfig(brokers kafka.Brokers, service service.Service, logger types.Logger) kafka.Config {
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

var Module = fx.Provide(
	GetKafkaBrokers,
	GetKafkaConfig,
	GetKafka,
)

func GetKafkaProtocol(kafka *kafka.Kafka) types.Protocol {
	return kafka
}

var ProtocolModule = fx.Provide(GetKafkaProtocol)
