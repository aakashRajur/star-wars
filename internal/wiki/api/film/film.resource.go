package film

import (
	"github.com/aakashRajur/star-wars/pkg/di/service-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL       = `/film/(?P<id>[0-9]+)`
	ParamFilmId   = `id`
	TypeFilmGet   = `FILM_GET`
	TypeFilmPatch = `FILM_PATCH`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypeFilmGet,
}

func Get() service_resource.ServiceResourceProvider {
	return service_resource.ServiceResourceProvider{
		Resource: ResourceGet,
	}
}

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypeFilmPatch,
}

func Patch() service_resource.ServiceResourceProvider {
	return service_resource.ServiceResourceProvider{
		Resource: ResourcePatch,
	}
}
