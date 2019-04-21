package planet

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/validations"
)

var PlanetValidation = map[string][]types.Validator{
	`name`:            {validations.ValidateString()},
	`rotation_period`: {validations.ValidateFloat()},
	`orbital_period`:  {validations.ValidateFloat()},
	`diameter`:        {validations.ValidateFloat()},
	`gravity`:         {validations.ValidateFloat()},
	`surface_water`:   {validations.ValidateFloat()},
	`population`:      {validations.ValidateInteger()},
	`description`:     {validations.ValidateString()},
	`climate`:         {validations.ValidateArray(validations.ValidateString())},
	`terrain`:         {validations.ValidateArray(validations.ValidateString())},
}
