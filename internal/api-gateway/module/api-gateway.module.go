package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
)

var ApiGatewayModule = fx.Options(
	module.FatalHandlerModule,
	module.EnvModule,
	module.LogModule,
	module.InstrumentationModule,
	module.ConsulModule,
	ResourceModule,
	module.HttpModule,
	module.HttpProtocolModule,
	module.AppModule,
)
