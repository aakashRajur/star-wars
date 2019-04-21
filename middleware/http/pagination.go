package http

import (
	"context"
	nativeHttp "net/http"
	"net/url"
	"strconv"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	PAGINATION = "PAGINATION"
)

func ParseQuery(query url.Values) (types.Pagination, error) {
	pagination := types.Pagination{}

	paginationIdString := query.Get(types.QueryPaginationId)
	if paginationIdString == "" {
		paginationIdString = "0"
	}
	paginationId, err := strconv.ParseInt(paginationIdString, 10, 0)
	if err != nil {
		return pagination, errors.NewBadRequest(err, "Invalid integer string provided for offset")
	}
	pagination.PaginationId = int64(paginationId)

	limitStr := query.Get(types.QueryLimit)
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.ParseInt(limitStr, 10, 0)
	if err != nil {
		return pagination, errors.NewBadRequest(err, "Invalid integer string provided for limit")
	}
	pagination.Limit = int64(limit)
	return pagination, nil
}

var Pagination http.Middleware = func(next http.HandleRequest) http.HandleRequest {
	return func(response http.Response, request *http.Request) {
		query := request.URL.Query()
		pagination, err := ParseQuery(query)
		if err != nil {
			response.Error(nativeHttp.StatusBadRequest, err)
		}

		withPagination := context.WithValue(request.Context(), PAGINATION, pagination)
		next(response, request.WithContext(withPagination))
	}
}
