package app

import (
	"github.com/aakashRajur/star-wars/internal/common/module/app"
	"github.com/aakashRajur/star-wars/internal/common/module/cache-strategy"
	"github.com/aakashRajur/star-wars/internal/common/module/consul"
	"github.com/aakashRajur/star-wars/internal/common/module/env"
	"github.com/aakashRajur/star-wars/internal/common/module/fatal"
	"github.com/aakashRajur/star-wars/internal/common/module/http"
	"github.com/aakashRajur/star-wars/internal/common/module/instrumentation"
	"github.com/aakashRajur/star-wars/internal/common/module/kafka"
	"github.com/aakashRajur/star-wars/internal/common/module/log"
	"github.com/aakashRajur/star-wars/internal/common/module/pg"
	"github.com/aakashRajur/star-wars/internal/common/module/redis"
	"github.com/aakashRajur/star-wars/internal/common/module/redis-pg"
	"github.com/aakashRajur/star-wars/internal/common/module/registree"
	"github.com/aakashRajur/star-wars/internal/topics/module/topics"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/module/healthceck"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/module/resource"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/module/service"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/module/subscription"
	"go.uber.org/fx"
)

var WikiKafkaModule = fx.Options(
	fatal.Module,
	env.Module,
	log.Module,
	instrumentation.Module,
	service.Module,
	pg.Module,
	redis.Module,
	redis_pg.Module,
	consul.Module,
	resource.Module,
	http.Module,
	healthceck.Module,
	kafka.Module,
	cache_strategy.Module,
	topics.Module,
	subscription.Module,
	kafka.ProtocolModule,
	app.Module,
	registree.Module,
)
