package validations

import (
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	ShapeError = `%s SHOULD BE AN OBJECT`
)

func ValidateShape(shape map[string][]types.Validator) types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}

		object, ok := value.(map[string]interface{})
		if !ok {
			return errors.Errorf(ShapeError, key)
		}

		for key, validators := range shape {
			value, ok := object[key]
			for _, validator := range validators {
				err := validator(key, value, ok)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}
