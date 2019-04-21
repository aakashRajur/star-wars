package validate

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func Validate(validators map[string][]types.Validator, body map[string]interface{}) error {
	tracker := make(map[string]bool, 1)
	errs := make(map[string]string, 1)
	allow := true

	for key, validators := range validators {
		tracker[key] = true

		for _, validator := range validators {
			value, ok := body[key]
			err := validator(key, value, ok)
			if err != nil {
				errs[key] = err.Error()
				allow = false
				break
			}
		}
	}

	for key := range body {
		_, ok := tracker[key]
		if ok {
			continue
		} else {
			errs[key] = fmt.Sprintf(`%s not allowed`, key)
			allow = false
		}
	}

	if !allow {
		errJson, err := json.Marshal(errs)
		if err != nil {
			return err
		} else {
			return errors.New(string(errJson))
		}
	}

	return nil
}

//noinspection GoNameStartsWithPackageName
func ValidateMapped(validators map[string][]types.Validator, body map[string]interface{}) map[string]string {
	tracker := make(map[string]bool, 1)
	errs := make(map[string]string, 1)
	allow := true

	for key, validators := range validators {
		tracker[key] = true

		for _, validator := range validators {
			value, ok := body[key]
			err := validator(key, value, ok)
			if err != nil {
				errs[key] = err.Error()
				allow = false
				break
			}
		}
	}

	for key := range body {
		_, ok := tracker[key]
		if ok {
			continue
		} else {
			errs[key] = fmt.Sprintf(`%s not allowed`, key)
			allow = false
		}
	}

	if !allow {
		return errs
	}

	return nil
}
