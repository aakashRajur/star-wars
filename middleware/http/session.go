package http

import (
	"context"
	nativeHttp "net/http"
	"time"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/util"
)

//noinspection GoSnakeCaseUsage
const (
	SESSION_COOKIE = `SESSION`
)

var Session http.Middleware = func(next http.HandleRequest) http.HandleRequest {
	return func(response http.Response, request *http.Request) {
		cookie, err := request.Cookie(SESSION_COOKIE)
		if err != nil || cookie.Value == `` {
			hex, err := util.SHA256()
			if err != nil {
				response.Error(nativeHttp.StatusInternalServerError, err)
				return
			}
			now := time.Now().UTC()
			expires := now.Add(8760 * time.Hour)

			cookie = &nativeHttp.Cookie{
				Name:     SESSION_COOKIE,
				Path:     `/`,
				Expires:  expires,
				HttpOnly: true,
				Value:    hex,
			}

			nativeHttp.SetCookie(response, cookie)
		}

		value := cookie.Value
		withSession := context.WithValue(request.Context(), SESSION_COOKIE, value)
		next(response, request.WithContext(withSession))
	}
}
