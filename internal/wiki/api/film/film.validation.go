package film

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/validations"
)

var ArgValidation = map[string][]types.Validator{
	`id`: {validations.ValidateRequired(), validations.ValidateIndex()},
}

var BodyValidation = map[string][]types.Validator{
	`title`:         {validations.ValidateString()},
	`episode`:       {validations.ValidateInteger()},
	`opening_crawl`: {validations.ValidateString()},
	`director`:      {validations.ValidateString()},
	`producer`:      {validations.ValidateString()},
	`release_date`:  {validations.ValidateTimestamp()},
	`description`:   {validations.ValidateString()},
	`planets`:       {validations.ValidateArray(validations.ValidateIndex())},
	`vehicles`:      {validations.ValidateArray(validations.ValidateIndex())},
	`species`:       {validations.ValidateArray(validations.ValidateIndex())},
	`characters`:    {validations.ValidateArray(validations.ValidateIndex())},
	`star_ships`:    {validations.ValidateArray(validations.ValidateIndex())},
}
