package film

import (
	"fmt"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"

	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func HttpResource(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, httpPrefix, film.HttpURL))
	resource.Get(GetHttpFilm(resolver, logger, tracker))
	resource.Patch(PatchHttpFilm(resolver, logger, tracker))

	return http_resource.HttpResourceProvider{Resource: resource}
}