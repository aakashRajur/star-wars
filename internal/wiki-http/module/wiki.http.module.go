package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
)

var WikiHttpModule = fx.Options(
	module.FatalHandlerModule,
	module.EnvModule,
	module.LogModule,
	module.InstrumentationModule,
	ServiceModule,
	module.PgModule,
	module.RedisModule,
	module.RedisPgModule,
	module.ConsulModule,
	module.HttpModule,
	module.HttpProtocolModule,
	module.CacheStrategyModule,
	ResourceModule,
	module.AppModule,
	module.RegistreeModule,
)
