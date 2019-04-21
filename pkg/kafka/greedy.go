package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/Shopify/sarama"

	"github.com/aakashRajur/star-wars/pkg/types"
)

var (
	timeoutPeriod = 5 * time.Second
)

var sortByTimestamp EventSorter = func(e1, e2 *Event) bool {
	return e1.Timestamp.Before(e2.Timestamp)
}

type greedy struct {
	ctx         context.Context
	until       time.Time
	logger      types.Logger
	events      map[string]Event
	topicDetail TopicDetail
	client      sarama.Client
	mux         sync.Mutex
}

func (g *greedy) ConsumePartition(partition int32) {
	logger := g.logger
	topicDetail := g.topicDetail

	consumer, err := sarama.NewConsumerFromClient(g.client)
	if err != nil {
		logger.Error(err)
		return
	}

	partitionConsumer, err := consumer.ConsumePartition(topicDetail.Name, partition, sarama.OffsetOldest)
	if err != nil {
		logger.Error(err)
		return
	}

	defer func() {
		err := partitionConsumer.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	ctx := g.ctx
	mux := g.mux
	until := g.until
	timeout, cancel := context.WithTimeout(context.Background(), timeoutPeriod)

	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				logger.Error(err)
				return
			}
		case <-timeout.Done():
			return
		case msg := <-partitionConsumer.Messages():
			if msg.Timestamp.After(until) {
				cancel()
				return
			}
			cancel()
			timeout, cancel = context.WithTimeout(context.Background(), timeoutPeriod)
			event, err := fromConsumerMessage(*msg, ctx)
			if err != nil {
				logger.Error(err)
				continue
			}
			mux.Lock()
			_, ok := g.events[event.Id]
			if !ok {
				g.events[event.Id] = event
			}
			mux.Unlock()
		}
	}
}

func (g *greedy) ConsumeAll() []Event {
	topicDetail := g.topicDetail
	partitions := int(topicDetail.NumPartitions)

	wg := sync.WaitGroup{}
	wg.Add(partitions)

	for i := 0; i < partitions; i++ {
		go func(partition int) {
			defer wg.Done()
			g.ConsumePartition(int32(partition))
		}(i)
	}

	wg.Wait()

	events := make([]Event, 0)
	for _, event := range g.events {
		events = append(events, event)
	}

	sortByTimestamp.Sort(events)

	return events
}

func NewGreedyInstance(ctx context.Context, client sarama.Client, logger types.Logger, topicDetail TopicDetail) *greedy {
	g := greedy{
		ctx:         ctx,
		until:       time.Now(),
		logger:      logger,
		events:      make(map[string]Event),
		topicDetail: topicDetail,
		client:      client,
		mux:         sync.Mutex{},
	}

	return &g
}
