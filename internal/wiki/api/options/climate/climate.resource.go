package climate

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpUrl        = `/options/climate`
	TypeClimateGet = `CLIMATE_GET`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpUrl,
	HttpVerb:   http.VerbGet,
	Type:       TypeClimateGet,
}
