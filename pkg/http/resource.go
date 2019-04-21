package http

import (
	"net/http"
	"regexp"
)

type Resource struct {
	Pattern  *regexp.Regexp
	Handlers map[string]HandleRequest
}

func (resource Resource) Get(handler WithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodGet] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Post(handler WithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodPost] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Put(handler WithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodPut] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Patch(handler WithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodPatch] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Delete(handler WithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodDelete] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Options(handler WithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodOptions] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	resolver := resource.Handlers
	var resolved = resolver[request.Method]
	if resolved == nil {
		http.NotFound(writer, request)
	} else {
		var req = Request{*request}

		resolved(
			Response{
				ResponseWriter: writer,
				compress:       req.CanGzip(),
			},
			&req,
		)
	}
}

func NewResource(pattern string) Resource {
	return Resource{
		Pattern:  regexp.MustCompile(pattern),
		Handlers: make(map[string]HandleRequest, 1),
	}
}
