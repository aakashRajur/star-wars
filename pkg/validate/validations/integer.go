package validations

import (
	"math"
	"reflect"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func ValidateInteger() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		elemName := reflect.TypeOf(value).Name()
		if elemName != `float64` {
			return errors.Errorf(`%v is not a valid type for %s, should be an integer`, value, key)
		}
		parsed := value.(float64)
		whole, fraction := math.Modf(parsed)
		if fraction > 0 || whole < 0 {
			return errors.Errorf(`%v is not a valid type for %s, should be an integer`, value, key)
		}
		return nil
	}
}
