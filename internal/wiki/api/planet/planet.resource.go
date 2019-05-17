package planet

import (
	"github.com/aakashRajur/star-wars/pkg/di/service-resource"
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

func Get() service_resource.ServiceResourceProvider {
	return service_resource.ServiceResourceProvider{
		Resource: ResourceGet,
	}
}

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypePlanetPatch,
}

func Patch() service_resource.ServiceResourceProvider {
	return service_resource.ServiceResourceProvider{
		Resource: ResourcePatch,
	}
}
