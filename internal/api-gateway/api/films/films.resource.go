package films

import (
	"fmt"

	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/observable"

	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func HttpResource(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, httpPrefix, films.HttpURL))
	resource.Get(GetHttpFilms(resolver, logger, tracker))

	return http_resource.HttpResourceProvider{Resource: resource}
}

func KafkaResource(resolver service.Resolver, kafkaInstance *kafka.Kafka, observable *observable.Observable, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics) http_resource.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, kafkaPrefix, films.HttpURL))
	resource.Get(GetKafkaFilms(resolver, kafkaInstance, observable, logger, tracker, definedTopics))

	return http_resource.HttpResourceProvider{Resource: resource}
}
