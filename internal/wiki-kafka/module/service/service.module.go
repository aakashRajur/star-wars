package service

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

//noinspection GoSnakeCaseUsage
const (
	SERVICE_NAME         = `SERVICE_NAME`
	HOSTNAME             = `CONTAINER_HOST_NAME`
	PORT                 = `HTTP_PORT`
	HEALTHCHECK_INTERVAL = `CONSUL_HEALTHCHECK_INTERVAL`
	INSTANCE_ID          = `INSTANCE_ID`
)

func GetService() service.Service {
	serviceName := env.GetString(SERVICE_NAME)
	hostName := env.GetString(HOSTNAME)
	port := env.GetInt(PORT)
	healthcheckInterval := env.GetString(HEALTHCHECK_INTERVAL)
	instanceId := env.GetString(INSTANCE_ID)

	healthcheck := service.Healthcheck{
		Scheme:   `http`,
		URL:      fmt.Sprintf(`http://%s:%d/stats`, hostName, port),
		HttpVerb: http.VerbGet,
		Interval: healthcheckInterval,
		Timeout:  `30s`,
	}

	return service.Service{
		Id:          instanceId,
		Name:        serviceName,
		Scheme:      `kafka`,
		Hostname:    hostName,
		Healthcheck: healthcheck,
	}
}

var Module = fx.Provide(GetService)
