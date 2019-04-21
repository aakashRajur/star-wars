package http

import (
	"context"
	"encoding/json"
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/pkg/http"
)

//noinspection GoSnakeCaseUsage
const (
	JSON_BODY = `JSON_BODY`
)

var JsonBodyParser http.Middleware = func(next http.HandleRequest) http.HandleRequest {
	return func(response http.Response, request *http.Request) {
		applicationType := request.Header.Get(http.ContentType)
		if applicationType != http.ContentTypeJSON {
			next(response, request)
			return
		}

		parsed := make(map[string]interface{}, 1)
		parser := json.NewDecoder(request.Body)

		err := parser.Decode(&parsed)
		if err != nil {
			response.Error(nativeHttp.StatusUnprocessableEntity, err)
			return
		}

		withParsedJson := context.WithValue(request.Context(), JSON_BODY, parsed)
		next(response, request.WithContext(withParsedJson))
	}
}
