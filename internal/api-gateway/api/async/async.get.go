package async

import (
	"context"
	"fmt"
	nativeHttp "net/http"
	"time"

	"github.com/juju/errors"

	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/observable"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetAsync(observable *observable.Observable, logger types.Logger, tracker types.TimeTracker) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		defer tracker(time.Now())

		nativeResponse := response.ResponseWriter

		ctx := request.Context()
		sessionValue := ctx.Value(middleware.SESSION_COOKIE)
		if sessionValue == nil {
			response.Error(
				nativeHttp.StatusInternalServerError,
				errors.New(`NO SESSION SET`),
			)
			return
		}
		session, ok := sessionValue.(string)
		if !ok {
			response.Error(
				nativeHttp.StatusInternalServerError,
				errors.New(`CORRUPT SESSION`),
			)
			return
		}

		header := response.Header()
		header.Set(http.ContentType, http.ContentTypeTextStream)
		header.Set(http.CacheControl, http.CacheControlNoCache)
		header.Set(http.Connection, http.ConnectionKeepAlive)

		err := observable.Register(session)
		if err != nil {
			logger.Error(err)
		}

		broker, err := observable.Broker(session)
		if err != nil {
			response.Error(
				nativeHttp.StatusInternalServerError,
				err,
			)
			return
		}

		flusher, ok := nativeResponse.(nativeHttp.Flusher)
		if !ok {
			response.Error(
				nativeHttp.StatusPreconditionFailed,
				errors.New(`SSE UNSUPPORTED BY CLIENT`),
			)
			return
		}

		_, _ = fmt.Fprintf(nativeResponse, ": handshake\n\n")
		flusher.Flush()

	ITERATOR:
		for {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil && err != context.Canceled {
					logger.Error(err)
				}
				err = observable.Unregister(session)
				if err != nil {
					logger.Error(err)
				}
				break ITERATOR
			case data := <-broker:
				if data == nil || len(data) < 1 {
					continue
				}
				_, err = fmt.Fprintf(nativeResponse, "data: %s\n\n", string(data))
				if err != nil {
					logger.Error(err)
				}
				flusher.Flush()
			}
		}
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.Session,
	)

	return http.HandlerWithMiddleware{
		Middlewares:   middlewares,
		HandleRequest: requestHandler,
	}
}
