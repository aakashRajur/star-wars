package normalizations

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NormalizeArray(elementNormalizer types.Normalizor) types.Normalizor {
	return func(key string, value interface{}) (interface{}, error) {
		if value == nil {
			return nil, nil
		}
		array, ok := value.([]interface{})
		compiled := ``

		if !ok {
			return nil, errors.Errorf(`%s is not a valid array`, key)
		}

		for i, value := range array {
			casted, err := elementNormalizer(key, value)
			if err != nil {
				return nil, errors.Errorf(`value at %d is not valid`, i)
			}
			if i > 0 {
				compiled += `, `
			}
			compiled += fmt.Sprintf(`%v`, casted)
		}

		return fmt.Sprintf(`{%s}`, compiled), nil
	}
}
