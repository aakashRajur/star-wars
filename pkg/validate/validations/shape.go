package validations

import (
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func ValidateShape(shape map[string][]types.Validator) types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}

		object, ok := value.(map[string]interface{})
		if !ok {
			return errors.Errorf(`%s should be an object`, key)
		}

		for k, v := range object {
			validators := shape[k]
			for _, validator := range validators {
				err := validator(k, v, true)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}
