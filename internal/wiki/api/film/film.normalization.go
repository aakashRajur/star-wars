package film

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
)

var FilmNormalization = map[string]types.Normalizor{
	`episode`:    normalizations.NormalizeInteger(),
	`planets`:    normalizations.NormalizeArray(normalizations.NormalizeInteger()),
	`vehicles`:   normalizations.NormalizeArray(normalizations.NormalizeInteger()),
	`species`:    normalizations.NormalizeArray(normalizations.NormalizeInteger()),
	`characters`: normalizations.NormalizeArray(normalizations.NormalizeInteger()),
	`star_ships`: normalizations.NormalizeArray(normalizations.NormalizeInteger()),
}
