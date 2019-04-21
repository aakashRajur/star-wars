package validations

import (
	"reflect"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
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
			return errors.Errorf(`%s should be a valid array`, key)
		}

		safeArray := value.([]interface{})
		if len(safeArray) == 0 {
			return nil
		}

		for _, each := range safeArray {
			err := elementValidator(key, each, true)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
