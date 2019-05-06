package module

import (
	"math"
	"strings"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetKafkaConfig(logger types.Logger, endpoint types.Endpoint) kafka.Config {
	brokers := strings.Split(env.GetString(`KAFKA_BROKERS`), `,`)
	maxPartitions := env.GetInt(`KAFKA_MAX_PARTITIONS`)
	maxReplicas := env.GetInt(`KAFKA_MAX_REPLICAS`)
	groupId := env.GetString(`KAFKA_GROUP_ID`)
	clientId := string(endpoint)

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

func GetKafka(handler types.FatalHandler, config kafka.Config, definedTopics kafka.DefinedTopics, subscriptions []*kafka.Subscription) *kafka.Kafka {
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
