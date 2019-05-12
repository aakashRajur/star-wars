package specie

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
)

var ArgNormalization = map[string]types.Normalizor{
	`id`: normalizations.NormalizeInteger(),
}

var BodyNormalization = map[string]types.Normalizor{
	`home_world`:  normalizations.NormalizeInteger(),
	`eye_colors`:  normalizations.NormalizeArray(normalizations.NormalizeString()),
	`hair_colors`: normalizations.NormalizeArray(normalizations.NormalizeString()),
	`skin_colors`: normalizations.NormalizeArray(normalizations.NormalizeString()),
}
