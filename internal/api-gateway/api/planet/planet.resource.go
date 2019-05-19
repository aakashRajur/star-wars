package planet

import (
	"fmt"

	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/observable"

	"github.com/aakashRajur/star-wars/internal/wiki/api/planet"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func HttpResource(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, httpPrefix, planet.HttpURL))
	resource.Get(GetHttpPlanet(resolver, logger, tracker))
	resource.Patch(PatchHttpPlanet(resolver, logger, tracker))

	return http_resource.HttpResourceProvider{Resource: resource}
}

func KafkaResource(resolver service.Resolver, kafkaInstance *kafka.Kafka, observable *observable.Observable, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) http_resource.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, kafkaPrefix, planet.HttpURL))
	resource.Get(GetKafkaPlanet(resolver, kafkaInstance, observable, logger, tracker, definedTopics, planet.ParamPlanetId))

	return http_resource.HttpResourceProvider{Resource: resource}
}
