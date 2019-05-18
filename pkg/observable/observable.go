package observable

import (
	"sync"

	"github.com/pkg/errors"
)

type Observable struct {
	brokers map[string]chan []byte
	mux     sync.Mutex
}

func (o *Observable) Register(id string) error {
	mux := o.mux
	mux.Lock()
	defer mux.Unlock()

	brokers := o.brokers
	_, ok := brokers[id]
	if ok {
		return errors.Errorf(`%s IS ALREADY REGISTERED`, id)
	}
	brokers[id] = make(chan []byte)

	return nil
}

func (o *Observable) Unregister(id string) error {
	mux := o.mux
	mux.Lock()
	defer mux.Unlock()

	brokers := o.brokers
	broker, ok := brokers[id]
	if !ok {
		return errors.Errorf(`%s WAS NEVER REGISTERED`, id)
	}

	close(broker)
	delete(brokers, id)

	return nil
}

func (o *Observable) Purge() {
	mux := o.mux
	mux.Lock()
	defer mux.Unlock()

	brokers := o.brokers
	for _, broker := range brokers {
		close(broker)
	}
}

func (o *Observable) SendData(id string, data []byte) error {
	mux := o.mux

	brokers := o.brokers

	mux.Lock()
	defer mux.Unlock()
	broker, ok := brokers[id]
	if !ok {
		return errors.Errorf(`%s WAS NEVER REGISTERED`, id)
	}

	broker <- data

	return nil
}

func (o *Observable) Broker(id string) (<-chan []byte, error) {
	mux := o.mux
	mux.Lock()
	defer mux.Unlock()

	brokers := o.brokers
	broker, ok := brokers[id]
	if !ok {
		return nil, errors.Errorf(`%s WAS NEVER REGISTERED`, id)
	}

	return broker, nil
}
