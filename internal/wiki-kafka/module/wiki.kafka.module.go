package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
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
	module.ConsulModule,
	module.KafkaModule,
	module.ResourceDefinitionModule,
	module.CacheStrategyModule,
	ResourceModule,
	HealthcheckModule,
	TopicsModule,
	SubscriptionModule,
	module.AppModule,
	module.ResourceRegistrationModule,
)
