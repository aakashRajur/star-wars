package app

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/api-gateway/module/kafka"
	"github.com/aakashRajur/star-wars/internal/api-gateway/module/resource"
	"github.com/aakashRajur/star-wars/internal/api-gateway/module/service"
	"github.com/aakashRajur/star-wars/internal/common/module/app"
	"github.com/aakashRajur/star-wars/internal/common/module/consul"
	"github.com/aakashRajur/star-wars/internal/common/module/env"
	"github.com/aakashRajur/star-wars/internal/common/module/fatal"
	"github.com/aakashRajur/star-wars/internal/common/module/http"
	"github.com/aakashRajur/star-wars/internal/common/module/instrumentation"
	"github.com/aakashRajur/star-wars/internal/common/module/log"
	"github.com/aakashRajur/star-wars/internal/common/module/observable"
	"github.com/aakashRajur/star-wars/internal/common/module/registree"
	"github.com/aakashRajur/star-wars/internal/common/module/subscription"
	"github.com/aakashRajur/star-wars/internal/topics/module/topics"
)

var ApiGatewayModule = fx.Options(
	fatal.Module,
	env.Module,
	log.Module,
	instrumentation.Module,
	service.Module,
	consul.Module,
	observable.Module,
	topics.Module,
	subscription.Module,
	kafka.Module,
	http.Module,
	http.ProtocolModule,
	resource.Module,
	app.Module,
	registree.Module,
)
