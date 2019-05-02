package http

import (
	"context"
	"net/http"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	PARAMS      = "PARAMS"
	HealthRoute = `/health`
	RootRoute   = `/`
)

type Router struct {
	*http.ServeMux
	Routes []Resource
	Health func() int
	Logger types.Logger
}

func (router *Router) DefineResource(resource Resource) *Router {
	router.Routes = append(router.Routes, resource)
	return router
}

func (router *Router) routeRequests(writer http.ResponseWriter, request *http.Request) bool {
	url := request.URL.Path

	if url == HealthRoute && router.Health != nil {
		status := router.Health()
		writer.WriteHeader(status)
		return true
	}

	params := make(map[string]string, 1)

	for _, each := range router.Routes {
		pattern := each.Pattern
		if pattern.MatchString(url) {
			matches := pattern.FindStringSubmatch(url)[1:]

			if len(matches) > 0 {
				keys := pattern.SubexpNames()[1:]
				for i := range keys {
					params[keys[i]] = matches[i]
				}
			}

			withParams := context.WithValue(request.Context(), PARAMS, params)
			each.ServeHTTP(
				writer,
				request.WithContext(withParams),
			)
			return true
		}
	}

	return false
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if router.routeRequests(writer, request) {
		return
	}

	http.NotFound(writer, request)
}

func NewRouter(logger types.Logger) *Router {
	router := &Router{
		Logger:   logger,
		ServeMux: http.NewServeMux(),
		Routes:   make([]Resource, 0),
	}

	return router
}
