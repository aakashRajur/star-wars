package character

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/validations"
)

var CharacterValidation = map[string][]types.Validator{
	`name`:        {validations.ValidateString()},
	`height`:      {validations.ValidateInteger()},
	`mass`:        {validations.ValidateFloat()},
	`hair_color`:  {validations.ValidateArray(validations.ValidateString())},
	`skin_color`:  {validations.ValidateArray(validations.ValidateString())},
	`eye_color`:   {validations.ValidateString()},
	`birth_year`:  {validations.ValidateString()},
	`gender`:      {validations.ValidateString()},
	`home_world`:  {validations.ValidateInteger()},
	`species`:     {validations.ValidateInteger()},
	`description`: {validations.ValidateString()},
	`vehicles`:    {validations.ValidateArray(validations.ValidateIndex())},
}
