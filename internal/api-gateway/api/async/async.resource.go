package async

import (
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/observable"
	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	HttpUrl = `/async`
)

func HttpResource(observable *observable.Observable, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(HttpUrl)
	resource.Get(GetAsync(observable, logger, tracker))

	return http_resource.HttpResourceProvider{Resource: resource}
}
