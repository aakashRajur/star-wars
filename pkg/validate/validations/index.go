package validations

import (
	"math"
	"reflect"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	IndexError = `%+v IS NOT A VALID INDEX FOR %s`
)

func ValidateIndex() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		elemName := reflect.TypeOf(value).Name()
		if elemName == `int` || elemName == `int32` || elemName == `int64` {
			return nil
		}
		if elemName != `float64` {
			return errors.Errorf(IndexError, value, key)
		}
		parsed := value.(float64)
		whole, fraction := math.Modf(parsed)
		if fraction > 0 || whole < 0 {
			return errors.Errorf(IndexError, parsed, key)
		}
		return nil
	}
}
