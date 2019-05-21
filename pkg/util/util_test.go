package util

import (
	"reflect"
	"testing"
	"time"
)

func TestMapStringFromInterfaces(t *testing.T) {
	type testCase struct {
		Args     []interface{}
		Expected []string
		Error    error
	}

	type complexType struct{}

	testCases := map[string]testCase{
		`ONLY_BASIC_TYPES`: {
			Args:     []interface{}{1, 2.5, true, []byte(`hello`), `world`},
			Expected: []string{`1`, `2.5`, `true`, `hello`, `world`},
			Error:    nil,
		},
		`WITH_STRUCTS`: {
			Args:     []interface{}{complexType{}, 1, 2.5, true, []byte(`hello`), `world`},
			Expected: []string{},
			Error:    ErrorUnsupportedType,
		},
	}

	for name, test := range testCases {
		var success = true

		expected := test.Expected
		expectedErr := test.Error
		got, err := MapStringFromInterfaces(test.Args)

		if !reflect.DeepEqual(got, expected) {
			success = false
			t.Errorf(`MapStringFromInterfaces() = %+v, want %+v`, got, expected)
		}
		if !reflect.DeepEqual(err, expectedErr) {
			success = false
			t.Errorf(`MapStringFromInterfaces() Error = %+v, want %+v`, err, expectedErr)
		}

		if success {
			t.Logf(`✔ %s`, name)
		}
	}
}

func TestGetType(t *testing.T) {
	type testCase struct {
		Args interface{}
		Type string
	}

	testCases := map[string]testCase{
		`INT`: {
			Args: int(10),
			Type: `int`,
		},
		`INT32`: {
			Args: int32(12),
			Type: `int32`,
		},
		`INT64`: {
			Args: int64(58),
			Type: `int64`,
		},
		`FLOAT32`: {
			Args: float32(5.0),
			Type: `float32`,
		},
		`FLOAT64`: {
			Args: float64(5.0),
			Type: `float64`,
		},
		`BOOL`: {
			Args: true,
			Type: `bool`,
		},
		`STRING`: {
			Args: `hello`,
			Type: `string`,
		},
	}

	for name, test := range testCases {
		got := GetType(test.Args)
		expected := test.Type
		if got != expected {
			t.Errorf(`GetType() = %+v, want %+v`, got, expected)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}

func TestDurationToString(t *testing.T) {
	type testCase struct {
		Args     time.Duration
		Expected string
	}

	testCases := map[string]testCase{
		`ONE_SECOND`: {
			Args:     time.Second,
			Expected: `1 second`,
		},
		`5_SECONDS`: {
			Args:     5 * time.Second,
			Expected: `5 seconds`,
		},
		`ONE_MINUTE`: {
			Args:     time.Minute,
			Expected: `1 minute`,
		},
		`10_MINUTES`: {
			Args:     10 * time.Minute,
			Expected: `10 minutes`,
		},
		`ONE_HOUR`: {
			Args:     time.Hour,
			Expected: `1 hour`,
		},
		`4_HOURS`: {
			Args:     4 * time.Hour,
			Expected: `4 hours`,
		},
		`ONE_DAYS`: {
			Args:     24 * time.Hour,
			Expected: `1 day`,
		},
		`3_DAYS`: {
			Args:     3 * 24 * time.Hour,
			Expected: `3 days`,
		},
		`30_MINUTES_20_SECONDS`: {
			Args:     30*time.Minute + 20*time.Second,
			Expected: `30 minutes 20 seconds`,
		},
		`8_HOURS_10_SECONDS`: {
			Args:     8*time.Hour + 10*time.Second,
			Expected: `8 hours 10 seconds`,
		},
		`10_DAYS_5_MINUTES`: {
			Args:     10*24*time.Hour + 5*time.Minute,
			Expected: `10 days 5 minutes`,
		},
		`2_DAYS_7_HOURS_18_MINUTES_48_SECONDS`: {
			Args:     2*24*time.Hour + 7*time.Hour + 18*time.Minute + 48*time.Second,
			Expected: `2 days 7 hours 18 minutes 48 seconds`,
		},
	}

	for name, test := range testCases {
		got := DurationToString(test.Args)
		expected := test.Expected
		if got != test.Expected {
			t.Errorf(`DurationToString() = %+v, want %+v`, got, expected)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}

func TestSHA256(t *testing.T) {
	hash, err := RandomSHA256()
	if err != nil {
		t.Errorf("RandomSHA256() = %+v, want %+v", err, nil)
	} else {
		t.Logf(`✔ %s`, hash)
	}
}
