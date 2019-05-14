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
	ServiceModule,
	module.ResourceModule,
	module.ConsulModule,
	module.HttpModule,
	module.HttpProtocolModule,
	ResourceModule,
	module.AppModule,
	module.RegistreeModule,
)
