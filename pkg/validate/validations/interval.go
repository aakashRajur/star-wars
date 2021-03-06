package validations

import (
	"encoding/json"
	"fmt"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	IntervalError = `%s is not allowed`
)

var intervalShapeValidator = ValidateShape(
	map[string][]types.Validator{
		`year`:  {ValidateInteger()},
		`month`: {ValidateInteger()},
		`day`:   {ValidateInteger()},
		`hour`:  {ValidateInteger()},
		`min`:   {ValidateInteger()},
		`sec`:   {ValidateInteger()},
	},
)

func ValidateInterval() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}

		err := intervalShapeValidator(key, value, exists)
		if err != nil {
			return err
		}

		allowedKeys := map[string]bool{
			`year`:  false,
			`month`: false,
			`day`:   false,
			`hour`:  false,
			`min`:   false,
			`sec`:   false,
		}
		interval := value.(map[string]interface{})
		errs := make(map[string]string, 0)

		for key := range interval {
			_, ok := allowedKeys[key]
			if !ok {
				errs[key] = fmt.Sprintf(IntervalError, key)
			}
		}

		if len(errs) > 0 {
			errJson, err := json.Marshal(errs)
			if err != nil {
				return err
			} else {
				return errors.New(string(errJson))
			}
		}

		return nil
	}
}
