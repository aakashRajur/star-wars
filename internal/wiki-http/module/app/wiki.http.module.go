package app

import (
	"github.com/aakashRajur/star-wars/internal/common/module/app"
	"github.com/aakashRajur/star-wars/internal/common/module/cache-strategy"
	"github.com/aakashRajur/star-wars/internal/common/module/consul"
	"github.com/aakashRajur/star-wars/internal/common/module/env"
	"github.com/aakashRajur/star-wars/internal/common/module/fatal"
	"github.com/aakashRajur/star-wars/internal/common/module/http"
	"github.com/aakashRajur/star-wars/internal/common/module/instrumentation"
	"github.com/aakashRajur/star-wars/internal/common/module/log"
	"github.com/aakashRajur/star-wars/internal/common/module/pg"
	"github.com/aakashRajur/star-wars/internal/common/module/redis"
	"github.com/aakashRajur/star-wars/internal/common/module/redis-pg"
	"github.com/aakashRajur/star-wars/internal/common/module/registree"
	"github.com/aakashRajur/star-wars/internal/wiki-http/module/resource"
	"github.com/aakashRajur/star-wars/internal/wiki-http/module/service"
	"go.uber.org/fx"
)

var WikiHttpModule = fx.Options(
	fatal.Module,
	env.Module,
	log.Module,
	instrumentation.Module,
	service.Module,
	pg.Module,
	redis.Module,
	redis_pg.Module,
	consul.Module,
	http.Module,
	http.ProtocolModule,
	cache_strategy.Module,
	resource.Module,
	app.Module,
	registree.Module,
)
