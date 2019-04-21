package normalizations

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NormalizeInterval() types.Normalizor {
	return func(key string, value interface{}) (interface{}, error) {
		if value == nil {
			return nil, nil
		}
		interval, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.Errorf(`%s is not a valid json for interval`, key)
		}

		compiled := ``
		unitCount := 0
		normalizor := NormalizeInteger()
		errs := make(map[string]string, 1)

		for key, value := range interval {
			val, err := normalizor(key, value)
			if err != nil {
				errs[key] = err.Error()
				continue
			}
			friendly := fmt.Sprintf(`%d %s`, val, key)
			if unitCount < 1 {
				compiled = friendly
			} else {
				compiled += fmt.Sprintf(` %s`, friendly)
			}
			unitCount += 1
		}

		if len(errs) > 0 {
			errJson, err := json.Marshal(errs)
			if err != nil {
				return nil, err
			} else {
				return nil, errors.New(string(errJson))
			}
		}

		return compiled, nil
	}
}
