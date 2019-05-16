package specie

import (
	"fmt"

	"github.com/aakashRajur/star-wars/internal/wiki/api/specie"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func HttpResource(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) di.HttpResourceProvider {
	resource := http.NewResource(fmt.Sprintf(`%s%s`, httpPrefix, specie.HttpURL))
	resource.Get(GetHttpSpecie(resolver, logger, tracker))

	return di.HttpResourceProvider{Resource: resource}
}
