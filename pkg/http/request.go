package http

import (
	"context"
	"net/http"
	"strings"
)

const (
	AcceptEncoding = `Accept-Encoding`
)

type Request struct {
	http.Request
}

func (request *Request) WithContext(ctx context.Context) *Request {
	original := request.Request.WithContext(ctx)
	return &(Request{*original})
}

func (request *Request) GetParams() map[string]interface{} {
	ctx := request.Context()
	return ctx.Value(PARAMS).(map[string]interface{})
}

func (request *Request) WithParams(params map[string]interface{}) *Request {
	withUpdatedParams := context.WithValue(
		request.Context(),
		PARAMS,
		params,
	)
	return request.WithContext(withUpdatedParams)
}

func (request *Request) CanGzip() bool {
	acceptedEncodings := request.Header.Get(AcceptEncoding)
	return strings.Contains(acceptedEncodings, ContentEncodingGzip)
}
