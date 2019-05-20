package skin_color

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpUrl        = `/options/skin-color`
	TypeClimateGet = `SKIN_COLOR_GET`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpUrl,
	HttpVerb:   http.VerbGet,
	Type:       TypeClimateGet,
}
