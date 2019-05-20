package vehicle

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL          = `/vehicle/(?P<id>[0-9]+)`
	ParamVehicleId   = `id`
	TypeVehicleGet   = `VEHICLE_GET`
	TypeVehiclePatch = `VEHICLE_PATCH`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypeVehicleGet,
}

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypeVehiclePatch,
}
