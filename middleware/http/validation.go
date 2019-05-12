package http

import (
	"context"
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/pkg/validate"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func ValidateArgs(validators map[string][]types.Validator, normalizors map[string]types.Normalizor) http.Middleware {
	return func(next http.HandleRequest) http.HandleRequest {
		return func(response http.Response, request *http.Request) {
			params := request.GetParams()

			err := validate.Validate(validators, params)
			if err != nil {
				response.ErrorString(nativeHttp.StatusUnprocessableEntity, http.ContentTypeJSON, err.Error())
			}

			if normalizors == nil {
				next(response, request)
				return
			}

			normalized, err := validate.Normalize(normalizors, params)
			if err != nil {
				response.ErrorString(nativeHttp.StatusUnprocessableEntity, http.ContentTypeJSON, err.Error())
				return
			}

			next(response, request.WithParams(normalized))
		}
	}
}

func ValidateBody(validators map[string][]types.Validator, normalizors map[string]types.Normalizor) http.Middleware {
	return func(next http.HandleRequest) http.HandleRequest {
		return func(response http.Response, request *http.Request) {
			ctx := request.Context()
			body := ctx.Value(JSON_BODY).(map[string]interface{})

			err := validate.Validate(validators, body)
			if err != nil {
				response.ErrorString(nativeHttp.StatusUnprocessableEntity, http.ContentTypeJSON, err.Error())
			}

			if normalizors == nil {
				next(response, request)
				return
			}

			normalized, err := validate.Normalize(normalizors, body)
			if err != nil {
				response.ErrorString(nativeHttp.StatusUnprocessableEntity, http.ContentTypeJSON, err.Error())
				return
			}

			normalizedBody := context.WithValue(request.Context(), JSON_BODY, normalized)
			next(response, request.WithContext(normalizedBody))
		}
	}
}
