package kafka

import (
	"github.com/Shopify/sarama"

	"github.com/aakashRajur/star-wars/pkg/types"
)

var (
	OffsetEarliest = sarama.OffsetOldest
	OffsetLatest   = sarama.OffsetNewest
)

type Config struct {
	Logger      types.Logger
	Brokers     []string
	GroupId     string
	ClientId    string
	Verbose     bool
	Replication int
	Partitions  int
}

func (config Config) ConsumerConfig(offset int64) *sarama.Config {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.V2_2_0_0
	consumerConfig.ClientID = config.ClientId
	consumerConfig.Consumer.Offsets.Initial = offset
	return consumerConfig
}

func (config Config) ProducerConfig() *sarama.Config {
	producerConfig := sarama.NewConfig()
	producerConfig.Version = sarama.V2_2_0_0
	producerConfig.ClientID = config.ClientId
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 3
	producerConfig.Producer.Return.Successes = true

	return producerConfig
}
