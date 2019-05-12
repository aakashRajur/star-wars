package planet

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL         = `/planet/(?P<id>[0-9]+)`
	ParamPlanetId   = `id`
	TypePlanetGet   = `PLANET_GET`
	TypePlanetPatch = `PLANET_PATCH`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypePlanetGet,
}

func Get() di.ServiceResourceProvider {
	return di.ServiceResourceProvider{
		Resource: ResourceGet,
	}
}

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypePlanetPatch,
}

func Patch() di.ServiceResourceProvider {
	return di.ServiceResourceProvider{
		Resource: ResourcePatch,
	}
}
