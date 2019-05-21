package validations

import (
	"fmt"
	"reflect"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	ArrayError = `%s SHOULD BE A VALID ARRAY`
)

func ValidateArray(elementValidator types.Validator) types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}

		rt := reflect.TypeOf(value)
		k := rt.Kind().String()
		kind := string(k)

		if kind != `array` && kind != `slice` {
			return errors.Errorf(ArrayError, key)
		}

		safeArray := value.([]interface{})
		if len(safeArray) == 0 {
			return nil
		}

		for i, each := range safeArray {
			err := elementValidator(fmt.Sprintf(`%s[%d]`, key, i), each, true)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
