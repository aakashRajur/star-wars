package character

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL            = `/character/(?P<id>[0-9]+)`
	ParamCharacterId   = `id`
	TypeCharacterGet   = `CHARACTER_GET`
	TypeCharacterPatch = `CHARACTER_PATCH`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypeCharacterGet,
}

var ResourcePatch = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbPatch,
	Type:       TypeCharacterPatch,
}
