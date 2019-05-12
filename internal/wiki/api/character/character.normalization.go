package character

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
)

var ArgNormalization = map[string]types.Normalizor{
	`id`: normalizations.NormalizeInteger(),
}

var BodyNormalization = map[string]types.Normalizor{
	`home_world`: normalizations.NormalizeInteger(),
	`species`:    normalizations.NormalizeInteger(),
	`hair_color`: normalizations.NormalizeArray(normalizations.NormalizeString()),
	`skin_color`: normalizations.NormalizeArray(normalizations.NormalizeString()),
	`vehicles`:   normalizations.NormalizeArray(normalizations.NormalizeInteger()),
}
