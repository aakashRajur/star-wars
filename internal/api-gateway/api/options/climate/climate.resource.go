package climate

import (
	"fmt"

	"github.com/aakashRajur/star-wars/internal/wiki/api/options/climate"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func HttpResource(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, httpPrefix, climate.HttpUrl))
	resource.Get(GetHttpClimate(resolver, logger, tracker))

	return http_resource.HttpResourceProvider{Resource: resource}
}
