package kafka

import (
	"go.uber.org/fx"

	kafkaModule "github.com/aakashRajur/star-wars/internal/common/module/kafka"
	"github.com/aakashRajur/star-wars/pkg/kafka"
)

func GetKafkaWithLifecycle(lifecycle fx.Lifecycle, kafkaInstance *kafka.Kafka) {
	lifecycle.Append(
		fx.Hook{
			OnStart: kafkaInstance.Start,
			OnStop:  kafkaInstance.Stop,
		},
	)
}

var Module = fx.Options(
	kafkaModule.Module,
	fx.Invoke(GetKafkaWithLifecycle),
)
