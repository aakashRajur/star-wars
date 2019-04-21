package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const ProtocolName string = `KAFKA`

type Kafka struct {
	Config        Config
	hook          Hook
	logger        types.Logger
	mux           sync.Mutex
	ctx           context.Context
	cancel        func()
	done          chan bool
	health        chan bool
	subscriptions map[string][]*Subscription
	client        sarama.Client
	consumer      sarama.ConsumerGroup
	producer      sarama.SyncProducer
}

func (kafka *Kafka) listening(topic string) (bool, error) {
	topics, err := kafka.client.Topics()
	if err != nil {
		return false, err
	}

	for _, each := range topics {
		if each == topic {
			return true, nil
		}
	}
	return false, nil
}

func (kafka *Kafka) controller() (*sarama.Broker, error) {
	return kafka.client.Controller()
}

func (kafka *Kafka) listTopic(topic string) (TopicDetail, error) {
	controller, err := kafka.controller()
	if err != nil {
		return TopicDetail{}, err
	}

	metaReq := &sarama.MetadataRequest{
		Version:                2,
		Topics:                 []string{topic},
		AllowAutoTopicCreation: false,
	}
	metaRes, err := controller.GetMetadata(metaReq)
	if err != nil {
		return TopicDetail{}, err
	}

	for _, meta := range metaRes.Topics {
		if meta.Name == topic {
			return FromTopicMetadata(meta), nil
		}
	}

	return TopicDetail{}, errors.Errorf(`topic: %s does not exist`, topic)
}

func (kafka *Kafka) listTopics() (map[string]TopicDetail, error) {
	controller, err := kafka.controller()
	if err != nil {
		return nil, err
	}

	metaReq := &sarama.MetadataRequest{}
	metaRes, err := controller.GetMetadata(metaReq)
	if err != nil {
		return nil, err
	}

	topicsDetailsMap := make(map[string]TopicDetail)

	var describeConfigsResources []*sarama.ConfigResource

	for _, topic := range metaRes.Topics {
		topicsDetailsMap[topic.Name] = FromTopicMetadata(topic)

		// we populate the resources we want to describe from the MetadataResponse
		topicResource := sarama.ConfigResource{
			Type: sarama.TopicResource,
			Name: topic.Name,
		}
		describeConfigsResources = append(describeConfigsResources, &topicResource)
	}

	// Send the DescribeConfigsRequest
	describeConfigsReq := &sarama.DescribeConfigsRequest{
		Resources: describeConfigsResources,
	}
	describeConfigsResp, err := controller.DescribeConfigs(describeConfigsReq)
	if err != nil {
		return nil, err
	}

	for _, resource := range describeConfigsResp.Resources {
		topicDetails := topicsDetailsMap[resource.Name]
		topicDetails.ConfigEntries = make(map[string]*string)

		for _, entry := range resource.Configs {
			// only include non-default non-sensitive config
			// (don't actually think topic config will ever be sensitive)
			if entry.Default || entry.Sensitive {
				continue
			}
			topicDetails.ConfigEntries[entry.Name] = &entry.Value
		}

		topicsDetailsMap[resource.Name] = topicDetails
	}

	return topicsDetailsMap, nil
}

func (kafka *Kafka) consumeMessages() {
	logger := kafka.Config.Logger
	failed := 0
	for {
		select {
		case <-kafka.ctx.Done():
			err := kafka.ctx.Err()
			if err != nil {
				logger.Error(err)
			}
			return
		default:
		}

		topics, err := kafka.client.Topics()
		if err != nil {
			logger.ErrorFields(
				err,
				map[string]interface{}{
					`msg`: `will wait and retry`,
				},
			)
			time.Sleep(1 * time.Second)
			continue
		}

		err = kafka.consumer.Consume(kafka.ctx, topics, kafka)
		if err != nil {
			logger.Error(err)
			failed += 1
		}
		if failed > 5 {
			failed = 0
			logger.InfoFields(
				map[string]interface{}{
					`reason`: fmt.Sprintf(`failed %d times`, failed),
					`delay`:  `1 second`,
					`err`:    err,
				},
			)
			time.Sleep(1 * time.Second)
		}
	}
}

func (kafka *Kafka) Subscribe(subscription *Subscription) error {
	topic := subscription.Topic

	kafka.mux.Lock()
	defer kafka.mux.Unlock()

	listening, err := kafka.listening(topic)
	if err != nil {
		return err
	}
	if !listening {
		return errors.Errorf(`topic: %s not initialized`, topic)
	}
	logger := kafka.Config.Logger
	logger.Info(fmt.Sprintf("subscribed to topic: %s", topic))

	subscribers, ok := kafka.subscriptions[topic]
	if !ok || len(subscribers) < 1 {
		subscribers = []*Subscription{subscription}
	} else {
		subscribers = append(subscribers, subscription)
	}
	kafka.subscriptions[topic] = subscribers

	return nil
}

func (kafka *Kafka) Unsubscribe(subscription *Subscription) error {
	topic := subscription.Topic

	kafka.mux.Lock()
	defer kafka.mux.Unlock()

	subscribers, ok := kafka.subscriptions[topic]
	if !ok {
		return nil
	}

	found := -1
	for index, each := range subscribers {
		if each == subscription {
			found = index
			break
		}
	}

	if found == -1 {
		return nil
	}

	subscribers = append(subscribers[:found], subscribers[found+1:]...)
	kafka.subscriptions[topic] = subscribers

	return nil
}

func (kafka *Kafka) SubscribeOnce(subscription *Subscription) error {
	var unsubscriber *Subscription = nil

	unsubscriber = &Subscription{
		Topic: subscription.Topic,
		Type:  subscription.Type,
		Id:    subscription.Id,
		Handler: func(event Event, k *Kafka) {
			logger := kafka.Config.Logger
			err := kafka.Unsubscribe(unsubscriber)
			if err != nil {
				logger.Error(err)
			}

			subscription.Handler(event, k)
		},
	}

	return kafka.Subscribe(unsubscriber)
}

func (kafka *Kafka) ConsumeFromBeginning(topic string) ([]Event, error) {
	topicDetail, err := kafka.listTopic(topic)
	if err != nil {
		return nil, err
	}

	ctx := kafka.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	g := NewGreedyInstance(
		ctx,
		kafka.client,
		kafka.logger,
		topicDetail,
	)

	events := g.ConsumeAll()

	return events, nil
}

func (kafka *Kafka) Emit(event Event) error {
	event.Source = kafka.Config.ClientId

	producer := kafka.producer
	message, err := event.ProducerMessage()

	if err != nil {
		return err
	}

	partition, offset, err := producer.SendMessage(&message)
	if err != nil {
		return err
	}

	kafka.logger.InfoFields(
		map[string]interface{}{
			`partition`: partition,
			`offset`:    offset,
			`id`:        event.Id,
			`topic`:     message.Topic,
		},
	)

	return nil
}

func (kafka *Kafka) CreateTopics(topics ...string) error {
	logger := kafka.Config.Logger
	refresh := make([]string, 0)
	availableTopics, err := kafka.listTopics()
	if err != nil {
		return err
	}

	config := kafka.Config
	topicDetails := make(map[string]*sarama.TopicDetail)

	for i, topic := range topics {
		if topic == `` {
			logger.Warn(fmt.Sprintf(`empty topic name found at %d`, i))
			continue
		}

		details, ok := availableTopics[topic]
		if ok {
			refresh = append(refresh, topic)
			logger.InfoFields(
				map[string]interface{}{
					`topic`:                  topic,
					`partitions`:             details.NumPartitions,
					`replication-factor`:     details.ReplicationFactor,
					`replication-assignment`: details.ReplicaAssignment,
				},
			)
			continue
		}

		logger.Info(fmt.Sprintf(`will create topic: %s`, topic))
		topicDetails[topic] = &sarama.TopicDetail{
			NumPartitions:     int32(config.Partitions),
			ReplicationFactor: int16(config.Replication),
		}
	}

	if len(topicDetails) < 1 {
		logger.Warn(`no new topics to create`)
		return kafka.client.RefreshMetadata(refresh...)
	} else {
		refresh = make([]string, 0)
	}

	controller, err := kafka.controller()
	if err != nil {
		return err
	}

	topicRequest := &sarama.CreateTopicsRequest{
		TopicDetails: topicDetails,
		ValidateOnly: false,
		Timeout:      5 * time.Second,
		Version:      2,
	}

	topicResponse, err := controller.CreateTopics(topicRequest)
	if topicResponse == nil {
		return errors.New(`no response while creating topics`)
	}
	topicErrs := topicResponse.TopicErrors

	errs := make(map[string]string)
	for topic, err := range topicErrs {
		if err.Err != sarama.ErrNoError && err.Err != sarama.ErrTopicAlreadyExists {
			errs[topic] = err.Error()
		}
	}

	if len(errs) > 0 {
		return errors.Errorf(`%v`, errs)
	}

	availableTopics, err = kafka.listTopics()
	if err != nil {
		return err
	}

	for topic, details := range availableTopics {
		refresh = append(refresh, topic)
		_, ok := topicDetails[topic]
		if !ok {
			continue
		}

		logger.InfoFields(
			map[string]interface{}{
				`topic`:                  topic,
				`partitions`:             details.NumPartitions,
				`replication-factor`:     details.ReplicationFactor,
				`replication-assignment`: details.ReplicaAssignment,
			},
		)
	}

	return kafka.client.RefreshMetadata(refresh...)
}

func (kafka *Kafka) Setup(sarama.ConsumerGroupSession) error {
	kafka.mux.Lock()
	defer kafka.mux.Unlock()
	kafka.health <- true

	onStart := kafka.hook.OnStart
	if onStart != nil {
		onStart(kafka)
	}

	return nil
}

func (kafka *Kafka) Cleanup(sarama.ConsumerGroupSession) error {
	kafka.mux.Lock()
	defer kafka.mux.Unlock()
	kafka.health <- false

	onStop := kafka.hook.OnStop
	if onStop != nil {
		onStop(kafka)
	}

	return nil
}

func (kafka *Kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		event, err := fromConsumerMessage(*message, kafka.ctx)
		if err != nil {
			kafka.logger.Error(err)
			continue
		}

		kafka.mux.Lock()
		subscribers, ok := kafka.subscriptions[event.Topic]
		kafka.mux.Unlock()
		if !ok {
			continue
		}

		kafka.logger.Info(event)

		for _, subscriber := range subscribers {
			if kafka.Config.ClientId != event.Source && subscriber.HasSubscribed(event) {
				go subscriber.Handler(event, kafka)
				session.MarkMessage(message, kafka.Config.ClientId)
			}
		}
	}

	return nil
}

func (kafka *Kafka) Name() string {
	return ProtocolName
}

func (kafka *Kafka) Start(ctx context.Context) error {
	logger := kafka.Config.Logger
	kafka.ctx, kafka.cancel = context.WithCancel(context.Background())

	go kafka.consumeMessages()
	logger.Info(`waiting for consumer to attach`)

	select {
	case <-kafka.health:
		logger.Info(`attached as consumer`)
		return nil
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return err
		}
		return nil
	}
}

func (kafka *Kafka) Stop(ctx context.Context) error {
	logger := kafka.Config.Logger
	err := make(chan error, 1)

	go func(err chan error) {
		logger := kafka.Config.Logger

		kafka.cancel()
		errs := make([]error, 0)

		consumer := kafka.consumer
		e := consumer.Close()
		if e != nil {
			errs = append(errs, e)
		} else {
			logger.Info(`successfully closed kafka consumer`)
		}

		producer := kafka.producer
		e = producer.Close()
		if e != nil {
			errs = append(errs, e)
		} else {
			logger.Info(`successfully closed kafka producer`)
		}

		client := kafka.client
		e = client.Close()
		if e != nil {
			errs = append(errs, e)
		} else {
			logger.Info(`successfully closed kafka client`)
		}

		if len(errs) > 0 {
			err <- errors.Errorf(`%v`, errs)
			return
		}

		close(err)
	}(err)

	select {
	case <-ctx.Done():
		logger.Error(`failed to close kafka consumer or producer, timedout`)
		return ctx.Err()
	case e, ok := <-err:
		if !ok {
			logger.Info(`kafka closed successfully`)
			return nil
		}
		return e
	}
}
