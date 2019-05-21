package validations

import (
	"math"
	"reflect"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	IntegerError = `%v IS NOT A VALID TYPE FOR %s, SHOULD BE AN INTEGER`
)

func ValidateInteger() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		elemName := reflect.TypeOf(value).Name()
		if elemName == `int` || elemName == `int32` || elemName == `int64` {
			return nil
		}
		if elemName != `float64` {
			return errors.Errorf(IntegerError, value, key)
		}
		parsed := value.(float64)
		whole, fraction := math.Modf(parsed)
		if fraction > 0 || whole < 0 {
			return errors.Errorf(IntegerError, value, key)
		}
		return nil
	}
}
