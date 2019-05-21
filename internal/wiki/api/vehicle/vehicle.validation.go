package vehicle

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/validations"
)

var ArgValidation = map[string][]types.Validator{
	`id`: {validations.ValidateRequired(), validations.ValidateIndex()},
}

var BodyValidation = map[string][]types.Validator{
	`name`:                  {validations.ValidateString()},
	`model`:                 {validations.ValidateString()},
	`manufacturer`:          {validations.ValidateString()},
	`cost_in_credits`:       {validations.ValidateInteger()},
	`size`:                  {validations.ValidateFloat()},
	`max_atmospheric_speed`: {validations.ValidateFloat()},
	`crew`:                  {validations.ValidateInteger()},
	`passengers`:            {validations.ValidateInteger()},
	`cargo_capacity`:        {validations.ValidateInteger()},
	`consumables`:           {validations.ValidateInterval()},
	`hyperdrive_rating`:     {validations.ValidateInteger()},
	`mglt`:                  {validations.ValidateInteger()},
	`starship_class`:        {validations.ValidateString()},
}
