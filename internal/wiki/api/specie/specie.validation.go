package specie

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/validations"
)

var ArgValidation = map[string][]types.Validator{
	`id`: {validations.Required(), validations.ValidateIndex()},
}

var BodyValidation = map[string][]types.Validator{
	`name`:             {validations.ValidateString()},
	`classification`:   {validations.ValidateString()},
	`average_height`:   {validations.ValidateFloat()},
	`average_lifespan`: {validations.ValidateFloat()},
	`home_world`:       {validations.ValidateInteger()},
	`spoken_language`:  {validations.ValidateString()},
	`description`:      {validations.ValidateString()},
	`eye_colors`:       {validations.ValidateArray(validations.ValidateString())},
	`hair_colors`:      {validations.ValidateArray(validations.ValidateString())},
	`skin_colors`:      {validations.ValidateArray(validations.ValidateString())},
}
