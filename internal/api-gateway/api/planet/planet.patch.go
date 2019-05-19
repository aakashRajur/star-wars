package planet

import (
	nativeHttp "net/http"
	"net/url"
	"strings"
	"time"

	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func PatchHttpPlanet(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http.HandlerWithMiddleware {
	proxy := http.NewProxy(
		func(url *url.URL, host *url.URL) *url.URL {
			url.Scheme = host.Scheme
			url.Path = strings.Replace(url.Path, httpPrefix, ``, 1)
			url.Host = host.Host
			return url
		},
	)

	requestHandler := func(response http.Response, request *http.Request) {
		defer tracker(time.Now())

		hosts, err := resolver.Resolve(downstreamHttp)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		proxy.UpdateHosts(hosts)
		proxy.HandleRequest(response, request)
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
