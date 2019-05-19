package observable

import "sync"

func NewInstance() *Observable {
	return &Observable{
		brokers: make(map[string]chan Payload),
		mux:     sync.Mutex{},
	}
}
