package validations

import (
	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

func ValidateFloat() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		if t := util.GetType(value); t != `float64` && t != `float32` {
			return errors.Errorf(`%v is not a valid type for %s, should be float`, value, key)
		}
		return nil
	}
}
