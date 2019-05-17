package main

import (
	"encoding/json"
	"fmt"

	"github.com/aakashRajur/star-wars/pkg/kafka"
)

/*
func setup() (string, []string) {
	typePtr := flag.String(`type`, `CONSUMER`, `Type of client one of {CONSUMER|PRODUCER}`)

	flag.Parse()

	connectionType := *typePtr
	if connectionType == `` {
		connectionType = `CONSUMER`
	} else if connectionType != `CONSUMER` && connectionType != `PRODUCER` && connectionType != `GREEDY` {
		panic(`type can be only CONSUMER or PRODUCER`)
	}
	env.LoadEnv()
	brokers := strings.Split(env.GetString(`KAFKA_BROKERS`), `,`)

	fmt.Println(`You're a ` + connectionType)
	fmt.Printf("KAFKA_BROKERS: %v\n", brokers)

	return connectionType, brokers
}

func getLogger() types.logger {
	return log.NewInstance(logrus.DebugLevel, &logrus.TextFormatter{})
}

func startProducer(wait chan bool, instance *kafka.Kafka, messageTopic string, logger types.logger) {
	text := make(chan string)
	go func(sender chan string) {
		fmt.Println(`Enter Messages:`)
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			sender <- text
		}
	}(text)

	for {
		select {
		case <-wait:
			fmt.Println(`CLOSING`)
			return
		case msg := <-text:
			id := fmt.Sprintf(`%x`, sha3.Sum256([]byte(msg)))
			event := kafka.Event{
				Topic:  messageTopic,
				Id:     id,
				Type:   `TEST_TYPE`,
				Source: `PRODUCER`,
				Args: map[string]interface{}{
					`random`: `data`,
				},
				Data: map[string]interface{}{
					`message`: string(msg),
				},
				Timestamp: time.Now().UTC(),
			}

			err := instance.Emit(event)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func startConsumer(wait chan bool, instance *kafka.Kafka, messageTopic string, _ types.logger) {
	err := instance.Subscribe(
		&kafka.Subscription{
			Topic: messageTopic,
			Handler: func(event kafka.Event, instance *kafka.Kafka) {
				fmt.Printf("%+v\n", event)
			},
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Subscribed to %s\n", messageTopic)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	ctx, done := context.WithCancel(context.Background())

	err = instance.Start(ctx)
	if err != nil {
		panic(err)
	}

	<-sigint
	done()

	<-wait
}

func startGreedy(_ chan bool, instance *kafka.Kafka, messageTopic string, logger types.logger) {
	logger.Info(fmt.Sprintf("WILL TRY CONSUMING UPTILL NOW"))
	events, err := instance.ConsumeFromBeginning(messageTopic)
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("CONSUMED %d MESSAGES", len(events)))
	for _, event := range events {
		fmt.Printf("%v\n", event)
	}
}

func main() {
	logger := getLogger()
	connectionType, brokers := setup()

	kafkaConfig := kafka.Config{
		logger:      logger,
		Brokers:     brokers,
		GroupId:     `SAMPLE_GROUP`,
		ClientId:    `TEST`,
		Partitions:  len(brokers),
		Replication: len(brokers),
	}
	instance, err := kafka.NewInstance(kafkaConfig, kafka.Hook{})
	if err != nil {
		panic(err)
	}

	messageTopic := `topic-example`
	err = instance.CreateTopics(messageTopic)
	if err != nil {
		panic(err)
	}

	wait := interrupt.NotifyOnInterrupt(instance.Stop, 10*time.Second, os.Interrupt)

	switch connectionType {
	case `PRODUCER`:
		startProducer(wait, instance, messageTopic, logger)
	case `CONSUMER`:
		startConsumer(wait, instance, messageTopic, logger)
	case `GREEDY`:
		startGreedy(wait, instance, messageTopic, logger)
	}
}

*/

func main() {
	e := kafka.Event{
		Topic:  `star-wars_wiki-request`,
		Type:   `CHARACTER_GET`,
		Id:     `asdf`,
		Source: `GATEWAY`,
		Args: map[string]interface{}{
			`id`: 1,
		},
	}

	marshaled, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", string(marshaled))
}
