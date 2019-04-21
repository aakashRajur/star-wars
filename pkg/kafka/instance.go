package kafka

import (
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

func NewInstance(config Config, hook Hook) (*Kafka, error) {
	logger := config.Logger
	if config.Verbose {
		sarama.Logger = log.New(logger.Out(), `[kafka]`, log.LstdFlags|log.LUTC)
	}

	consumerConfig := config.ConsumerConfig(OffsetLatest)
	err := consumerConfig.Validate()
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewClient(config.Brokers, consumerConfig)
	if err != nil {
		return nil, err
	}
	failClientClose := func() {
		err := client.Close()
		if err != nil {
			logger.Error(err)
		}
	}

	consumer, err := sarama.NewConsumerGroupFromClient(
		config.GroupId,
		client,
	)
	if err != nil {
		defer failClientClose()
		return nil, err
	}
	failConsumerClose := func() {
		err := consumer.Close()
		if err != nil {
			logger.Error(err)
		}
	}

	producerConfig := config.ProducerConfig()
	err = producerConfig.Validate()
	if err != nil {
		defer failConsumerClose()
		defer failClientClose()
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, producerConfig)
	if err != nil {
		defer failConsumerClose()
		defer failClientClose()
		return nil, err
	}

	kafka := Kafka{
		Config:        config,
		hook:          hook,
		logger:        logger,
		mux:           sync.Mutex{},
		health:        make(chan bool),
		subscriptions: make(map[string][]*Subscription),
		client:        client,
		consumer:      consumer,
		producer:      producer,
	}

	return &kafka, nil
}
