package characters

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
)

const (
	HttpURL           = `/characters`
	TypeCharactersGet = `CHARACTERS_GET`
)

var ResourceGet = service.Resource{
	ApiPattern: HttpURL,
	HttpVerb:   http.VerbGet,
	Type:       TypeCharactersGet,
}

func Get() di.ServiceResourceProvider {
	return di.ServiceResourceProvider{
		Resource: ResourceGet,
	}
}
