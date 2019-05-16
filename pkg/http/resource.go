package http

import (
	"fmt"
	"net/http"
	"regexp"
)

type Resource struct {
	Pattern  *regexp.Regexp
	Handlers map[string]HandleRequest
}

func (resource Resource) Get(handler HandlerWithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodGet] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Post(handler HandlerWithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodPost] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Put(handler HandlerWithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodPut] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Patch(handler HandlerWithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodPatch] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Delete(handler HandlerWithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodDelete] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) Options(handler HandlerWithMiddleware) Resource {
	resolver := resource.Handlers
	resolver[http.MethodOptions] = handler.GetHTTPHandler()
	return resource
}

func (resource Resource) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	resolver := resource.Handlers
	resolved, ok := resolver[request.Method]
	if !ok {
		http.NotFound(writer, request)
		return
	}
	req := Request{*request}

	resolved(
		Response{
			ResponseWriter: writer,
			compress:       req.CanGzip(),
		},
		&req,
	)
}

func NewResource(pattern string) Resource {
	return Resource{
		Pattern:  regexp.MustCompile(fmt.Sprintf(`^%s$`, pattern)),
		Handlers: make(map[string]HandleRequest),
	}
}
