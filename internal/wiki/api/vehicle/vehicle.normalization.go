package vehicle

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
)

var ArgNormalization = map[string]types.Normalizor{
	`id`: normalizations.NormalizeInteger(),
}

var BodyNormalization = map[string]types.Normalizor{
	`cost_in_credits`:   normalizations.NormalizeInteger(),
	`crew`:              normalizations.NormalizeInteger(),
	`passengers`:        normalizations.NormalizeInteger(),
	`cargo_capacity`:    normalizations.NormalizeInteger(),
	`consumables`:       normalizations.NormalizeInterval(),
	`hyperdrive_rating`: normalizations.NormalizeInteger(),
	`mglt`:              normalizations.NormalizeInteger(),
}
