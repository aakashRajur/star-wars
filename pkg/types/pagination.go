package types

import (
	"encoding/json"
)

const (
	QueryTotalCount   = "total_count"
	QueryPaginationId = "pagination_id"
	QueryLimit        = "limit"
	PAGINATION        = "PAGINATION"
)

type Pagination struct {
	PaginationId    int64 `json:"pagination_id"`
	TotalCount      int64 `json:"total_count"`
	Limit           int64 `json:"limit"`
	LowestRecordId  int64 `json:"lowest_record_id"`
	HighestRecordId int64 `json:"highest_record_id"`
}

func PaginationFromString(str string) (Pagination, error) {
	var pagination Pagination
	err := json.Unmarshal([]byte(str), &pagination)
	return pagination, err
}
