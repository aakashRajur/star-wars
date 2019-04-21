package kafka

type HandleEvent func(event Event, instance *Kafka)

type Middleware func(handler HandleEvent) HandleEvent

func ChainMiddlewares(middlewares ...Middleware) Middleware {
	return func(handler HandleEvent) HandleEvent {
		next := handler
		for i := len(middlewares) - 1; i > -1; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
