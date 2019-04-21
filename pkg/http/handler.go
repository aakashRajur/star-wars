package http

type HandleRequest func(response Response, request *Request)

type Middleware func(request HandleRequest) HandleRequest

func ChainMiddlewares(middlewares ...Middleware) Middleware {
	return func(request HandleRequest) HandleRequest {
		var next = request
		for i := len(middlewares) - 1; i > -1; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

type WithMiddleware struct {
	Middlewares   Middleware
	HandleRequest HandleRequest
}

func (handler WithMiddleware) GetHTTPHandler() HandleRequest {
	if handler.Middlewares != nil {
		return handler.Middlewares(handler.HandleRequest)
	} else {
		return handler.HandleRequest
	}
}
