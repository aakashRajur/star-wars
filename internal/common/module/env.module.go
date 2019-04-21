package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/args"
	"github.com/aakashRajur/star-wars/pkg/env"
)

var EnvModule = fx.Options(
	fx.Invoke(args.LoadArgs),
	fx.Invoke(env.LoadEnv),
)
