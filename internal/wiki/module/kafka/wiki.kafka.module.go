package kafka

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
)

var WikiKafkaModule = fx.Options(
	module.FatalHandlerModule,
	module.EnvModule,
	module.LogModule,
	module.InstrumentationModule,
	module.PgModule,
	module.RedisModule,
	module.RedisPgModule,
	module.KafkaModule,
	module.CacheStrategyModule,
	ResourceDefinitionModule,
	HookModule,
	TopicsModule,
	SubscriptionModule,
	module.AppModule,
)
