package validations

import (
	"math"
	"reflect"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func ValidateIndex() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		elemName := reflect.TypeOf(value).Name()
		if elemName != `float64` {
			return errors.Errorf(`%v is not a valid index for %s`, value, key)
		}
		parsed := value.(float64)
		whole, fraction := math.Modf(parsed)
		if fraction > 0 || whole < 0 {
			return errors.Errorf(`%d is not a valid index for %s`, parsed, key)
		}
		return nil
	}
}
