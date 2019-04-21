package kafka

type Subscription struct {
	Topic   string
	Type    string
	Id      string
	Handler HandleEvent
}

func (subscription *Subscription) HasSubscribed(event Event) bool {
	if subscription.Topic != `` && subscription.Topic != event.Topic {
		return false
	} else if subscription.Type != `` && subscription.Type != event.Type {
		return false
	} else if subscription.Id != `` && subscription.Id != event.Id {
		return false
	}

	return true
}

func WithMiddleware(subscription *Subscription, middleware Middleware) *Subscription {
	subscription.Handler = middleware(subscription.Handler)
	return subscription
}
