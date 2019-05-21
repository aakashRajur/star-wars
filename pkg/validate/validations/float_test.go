package validations

import (
	"github.com/juju/errors"
	"testing"
)

func TestValidateFloat(t *testing.T) {
	type arg struct {
		Key    string
		Value  interface{}
		Exists bool
	}
	type testCase struct {
		Args     arg
		Expected error
	}

	testCases := map[string]testCase{
		`FLOAT`: {
			Args: arg{
				Key:    `key1`,
				Value:  10.5,
				Exists: true,
			},
			Expected: nil,
		},
		`INT`: {
			Args: arg{
				Key:    `key1`,
				Value:  16,
				Exists: true,
			},
			Expected: errors.Errorf(FloatError, `16`, `key1`),
		},
		`STRING`: {
			Args: arg{
				Key:    `key2`,
				Value:  `445.98`,
				Exists: true,
			},
			Expected: errors.Errorf(FloatError, `445.98`, `key2`),
		},
		`BOOL`: {
			Args: arg{
				Key:    `key3`,
				Value:  true,
				Exists: true,
			},
			Expected: errors.Errorf(FloatError, true, `key3`),
		},
		`EXISTS`: {
			Args: arg{
				Key:    `key4`,
				Value:  nil,
				Exists: false,
			},
			Expected: nil,
		},
	}

	validator := ValidateFloat()

	for name, test := range testCases {
		args := test.Args
		got := validator(args.Key, args.Value, args.Exists)
		expected := test.Expected
		if (got != nil && expected == nil) ||
			(got == nil && expected != nil) ||
			(got != nil && expected != nil && got.Error() != expected.Error()) {
			t.Errorf("ValidateFloat() = %+v, want %+v", got, expected)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
