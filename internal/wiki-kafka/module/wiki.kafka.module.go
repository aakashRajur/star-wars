package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
	topics "github.com/aakashRajur/star-wars/internal/topics/module"
)

var WikiKafkaModule = fx.Options(
	module.FatalHandlerModule,
	module.EnvModule,
	module.EndpointModule,
	module.LogModule,
	module.InstrumentationModule,
	module.PgModule,
	module.RedisModule,
	module.RedisPgModule,
	module.KafkaModule,
	module.ResourceDefinitionModule,
	module.CacheStrategyModule,
	topics.TopicsModule,
	SubscriptionModule,
	module.KafkaProtocolModule,
	module.AppModule,
)
