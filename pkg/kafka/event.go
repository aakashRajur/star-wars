package kafka

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/Shopify/sarama"
	"golang.org/x/crypto/sha3"
)

type Event struct {
	Topic     string                 `json:"topic" xml:"topic"`
	Type      string                 `json:"type" xml:"type"`
	Id        string                 `json:"id" xml:"id"`
	Source    string                 `json:"source" xml:"source"`
	Args      map[string]interface{} `json:"args" xml:"args"`
	Data      interface{}            `json:"data" xml:"data"`
	Error     map[string]string      `json:"error" xml:"error"`
	Timestamp time.Time              `json:"timestamp" xml:"timestamp"`
	Ctx       context.Context        `json:"-"`
}

func (event Event) ProducerMessage() (sarama.ProducerMessage, error) {
	event.Timestamp = time.Now()
	marshaled, err := json.Marshal(event)
	if err != nil {
		return sarama.ProducerMessage{}, err
	}

	k := sha3.Sum256([]byte(time.Now().UTC().String()))
	message := sarama.ProducerMessage{
		Topic:     event.Topic,
		Key:       sarama.ByteEncoder(k[:]),
		Value:     sarama.ByteEncoder(marshaled),
		Timestamp: event.Timestamp,
	}

	return message, nil
}

func fromConsumerMessage(message sarama.ConsumerMessage, ctx context.Context) (Event, error) {
	value := message.Value
	var event Event

	err := json.Unmarshal(value, &event)
	if err != nil {
		return Event{}, err
	}
	event.Ctx = ctx

	return event, nil
}

type EventSorter func(e1, e2 *Event) bool

func (es EventSorter) Sort(events []Event) {
	esh := &eventSorterHelper{
		events: events,
		by:     es,
	}
	sort.Sort(esh)
}

type eventSorterHelper struct {
	events []Event
	by     EventSorter
}

func (es *eventSorterHelper) Less(i, j int) bool {
	return es.by(&es.events[i], &es.events[j])
}

func (es *eventSorterHelper) Swap(i, j int) {
	es.events[i], es.events[j] = es.events[j], es.events[i]
}

func (es *eventSorterHelper) Len() int {
	return len(es.events)
}
