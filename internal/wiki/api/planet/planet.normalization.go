package planet

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
)

var ArgNormalization = map[string]types.Normalizor{
	`id`: normalizations.NormalizeInteger(),
}

var BodyNormalization = map[string]types.Normalizor{
	`population`: normalizations.NormalizeInteger(),
	`climate`:    normalizations.NormalizeArray(normalizations.NormalizeString()),
	`terrain`:    normalizations.NormalizeArray(normalizations.NormalizeString()),
}
