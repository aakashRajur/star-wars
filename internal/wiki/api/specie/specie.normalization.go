package specie

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
)

var SpecieNormalization = map[string]types.Normalizor{
	`home_world`:  normalizations.NormalizeInteger(),
	`eye_colors`:  normalizations.NormalizeArray(normalizations.NormalizeString()),
	`hair_colors`: normalizations.NormalizeArray(normalizations.NormalizeString()),
	`skin_colors`: normalizations.NormalizeArray(normalizations.NormalizeString()),
}
