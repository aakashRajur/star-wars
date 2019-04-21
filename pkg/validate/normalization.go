package validate

import (
	"encoding/json"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func Normalize(normalizors map[string]types.Normalizor, body map[string]interface{}) (map[string]interface{}, error) {
	normalized := make(map[string]interface{}, 1)
	errs := make(map[string]string)
	allow := true

	for key, value := range body {
		processed := value
		normalizor, ok := normalizors[key]
		if ok {
			val, err := normalizor(key, value)
			if err != nil {
				errs[key] = err.Error()
				allow = false
				continue
			}
			processed = val
		}
		normalized[key] = processed
	}

	if !allow {
		errJson, err := json.Marshal(errs)
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New(string(errJson))
		}
	}

	return normalized, nil
}

func NormalizeMapped(normalizors map[string]types.Normalizor, body map[string]interface{}) (map[string]interface{}, map[string]string) {
	normalized := make(map[string]interface{}, 1)
	errs := make(map[string]string)
	allow := true

	for key, value := range body {
		processed := value
		normalizor, ok := normalizors[key]
		if ok {
			val, err := normalizor(key, value)
			if err != nil {
				errs[key] = err.Error()
				allow = false
				continue
			}
			processed = val
		}
		normalized[key] = processed
	}

	if !allow {
		return nil, errs
	}

	return normalized, nil
}
