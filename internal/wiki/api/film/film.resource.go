package film

import (
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

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypeFilmPatch,
}
