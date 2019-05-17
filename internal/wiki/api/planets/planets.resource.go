package planets

import (
	"github.com/aakashRajur/star-wars/pkg/di/service-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL        = `/planets`
	TypePlanetsGet = `PLANETS_GET`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypePlanetsGet,
}

func Get() service_resource.ServiceResourceProvider {
	return service_resource.ServiceResourceProvider{
		Resource: ResourceGet,
	}
}
