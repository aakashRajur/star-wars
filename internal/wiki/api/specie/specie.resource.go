package specie

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL         = `/specie/(?P<id>[0-9]+)`
	ParamSpecieId   = `id`
	TypeSpecieGet   = `SPECIE_GET`
	TypeSpeciePatch = `SPECIE_PATCH`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypeSpecieGet,
}

func Get() di.ServiceResourceProvider {
	return di.ServiceResourceProvider{
		Resource: ResourceGet,
	}
}

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypeSpeciePatch,
}

func Patch() di.ServiceResourceProvider {
	return di.ServiceResourceProvider{
		Resource: ResourcePatch,
	}
}
