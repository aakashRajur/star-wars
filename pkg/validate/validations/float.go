package validations

import (
	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

const (
	FloatError = `%+v IS NOT A VALID TYPE FOR %s, SHOULD BE FLOAT`
)

func ValidateFloat() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		if t := util.GetType(value); t != `float64` && t != `float32` {
			return errors.Errorf(FloatError, value, key)
		}
		return nil
	}
}
