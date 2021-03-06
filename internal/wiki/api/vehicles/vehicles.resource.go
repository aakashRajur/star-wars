package vehicles

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL         = `/vehicles`
	TypeVehiclesGet = `VEHICLES_GET`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypeVehiclesGet,
}
