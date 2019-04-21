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

func (request *Request) GetParams() map[string]string {
	ctx := request.Context()
	return ctx.Value(PARAMS).(map[string]string)
}

func (request *Request) CanGzip() bool {
	acceptedEncodings := request.Header.Get(AcceptEncoding)
	return strings.Contains(acceptedEncodings, ContentEncodingGzip)
}
