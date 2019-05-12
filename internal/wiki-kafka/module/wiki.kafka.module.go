package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
	topics "github.com/aakashRajur/star-wars/internal/topics/module"
)

var WikiKafkaModule = fx.Options(
	module.FatalHandlerModule,
	module.EnvModule,
	module.LogModule,
	module.InstrumentationModule,
	ServiceModule,
	module.PgModule,
	module.RedisModule,
	module.RedisPgModule,
	module.ConsulModule,
	ResourceModule,
	module.HttpModule,
	HealthcheckModule,
	module.KafkaModule,
	module.CacheStrategyModule,
	topics.TopicsModule,
	SubscriptionModule,
	module.KafkaProtocolModule,
	module.AppModule,
	module.RegistreeModule,
)
