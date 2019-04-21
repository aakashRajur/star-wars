package kafka

type Hook struct {
	OnStart func(*Kafka)
	OnStop func(*Kafka)
}
