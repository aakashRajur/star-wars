package http

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/aakashRajur/star-wars/pkg/http"
	nativeHttp "net/http"
	"time"
)

//noinspection GoSnakeCaseUsage
const (
	SESSION_COOKIE = `SESSION`
)

var Session http.Middleware = func(next http.HandleRequest) http.HandleRequest {
	return func(response http.Response, request *http.Request) {
		cookie, err := request.Cookie(SESSION_COOKIE)
		if err != nil || cookie.Value == `` {
			now := time.Now().UTC()
			marshaled, err := now.MarshalBinary()
			if err != nil {
				response.Error(nativeHttp.StatusInternalServerError, err)
				return
			}
			sum := sha256.Sum256(marshaled)
			hex := fmt.Sprintf(`%x`, sum)

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
