package env

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/args"
	"github.com/aakashRajur/star-wars/pkg/env"
)

var Module = fx.Options(
	fx.Invoke(args.LoadArgs),
	fx.Invoke(env.LoadEnv),
)
