package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPaginationFromString(t *testing.T) {
	type testCase struct {
		Args     string
		Expected Pagination
		Error    error
	}

	testCases := map[string]testCase{
		`EMPTY_OBJECT`: {
			Args:     `{}`,
			Expected: Pagination{},
			Error:    nil,
		},
		`ONLY_PAGINATION_ID`: {
			Args: `{"pagination_id": 1}`,
			Expected: Pagination{
				PaginationId: 1,
			},
			Error: nil,
		},
		`ONLY_TOTAL_COUNT`: {
			Args: `{"total_count": 10}`,
			Expected: Pagination{
				TotalCount: 10,
			},
			Error: nil,
		},
		`ONLY_LIMIT`: {
			Args: `{"limit": 20}`,
			Expected: Pagination{
				Limit: 20,
			},
			Error: nil,
		},
		`ONLY_LOWEST_RECORD_ID`: {
			Args: `{"lowest_record_id": 15}`,
			Expected: Pagination{
				LowestRecordId: 15,
			},
			Error: nil,
		},
		`ONLY_HIGHEST_RECORD_ID`: {
			Args: `{"highest_record_id": 7}`,
			Expected: Pagination{
				HighestRecordId: 7,
			},
			Error: nil,
		},
		`PARTIAL`: {
			Args: `{"pagination_id": 5, "limit": 20}`,
			Expected: Pagination{
				PaginationId: 5,
				Limit:        20,
			},
			Error: nil,
		},
		`ARRAY`: {
			Args:     `[]`,
			Expected: Pagination{},
			Error: &json.UnmarshalTypeError{
				Value:  `array`,
				Type:   reflect.TypeOf(Pagination{}),
				Offset: 1,
				Struct: "",
				Field:  "",
			},
		},
	}

	for name, test := range testCases {
		success := true
		got, err := PaginationFromString(test.Args)
		expected := test.Expected
		expectedErr := test.Error

		if !reflect.DeepEqual(err, expectedErr) {
			success = false
			t.Errorf("PaginationFromString() error = %v, wantErr %v", err, expectedErr)
		}
		if !reflect.DeepEqual(got, expected) {
			success = false
			t.Errorf("PaginationFromString() = %+v, want %+v", got, expected)
		}

		if success {
			t.Logf(`✔ %s`, name)
		}
	}
}
