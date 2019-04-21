package http

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
)

var WikiHttpModule = fx.Options(
	module.FatalHandlerModule,
	module.EnvModule,
	module.LogModule,
	module.InstrumentationModule,
	module.PgModule,
	module.RedisModule,
	module.RedisPgModule,
	module.HttpModule,
	module.CacheStrategyModule,
	ResourceModule,
	module.AppModule,
)
