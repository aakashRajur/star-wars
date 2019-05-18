package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	nativeHttp "net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	keyId                        = `id`
	keyName                      = `name`
	keyAddress                   = `address`
	keyPort                      = `port`
	keyHealthcheck               = `check`
	keyHealthcheckSkipTLS        = `tls_skip_verify`
	keyHealthcheckMethod         = `method`
	keyHealthcheckInterval       = `interval`
	keyHealthcheckTimeout        = `timeout`
	consulServiceRegisterPath    = `/v1/agent/service/register`
	consulServiceUnregisterParam = `service_id`
	consulServiceUnregisterPath  = `/v1/agent/service/deregister/:service_id`
	consulWatchServicesPath      = `/v1/catalog/services`
	consulBlockingIndexKey       = `index`
	consulBlockingWaitKey        = `wait`
	consulBlockingWaitValue      = `5m`
	consulBlockingHeaderKey      = `X-Consul-Index`
	consulGetServicesPath        = `/v1/agent/services`
	resolveRetry                 = 5
	resolveWait                  = 5 * time.Second
)

func registerService(consul *Consul, definition service.Service) error {
	config := consul.config
	body := make(map[string]interface{})

	body[keyId] = definition.Id
	body[keyName] = definition.Name
	body[keyAddress] = fmt.Sprintf(
		`%s://%s:%d`,
		definition.Scheme,
		definition.Hostname,
		definition.Port,
	)
	body[keyPort] = definition.Port
	healthcheck := definition.Healthcheck
	body[keyHealthcheck] = map[string]interface{}{
		healthcheck.Scheme:     healthcheck.URL,
		keyHealthcheckSkipTLS:  healthcheck.SkipTLS,
		keyHealthcheckMethod:   healthcheck.HttpVerb,
		keyHealthcheckInterval: healthcheck.Interval,
		keyHealthcheckTimeout:  healthcheck.Timeout,
	}

	url := config.Url()
	url.Path = consulServiceRegisterPath

	requestConfig := http.RequestConfig{
		Verb:    http.VerbPut,
		Url:     url,
		Headers: nil,
		Params:  nil,
		Body:    body,
	}
	client := consul.client

	request, err := http.NewRequest(requestConfig)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO REGISTER SERVICE: %s`, response.Status)
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func unregisterService(consul *Consul, definition service.Service) error {
	config := consul.config
	url := config.Url()
	url.Path = consulServiceUnregisterPath

	requestConfig := http.RequestConfig{
		Verb:    `PUT`,
		Url:     url,
		Headers: nil,
		Params: map[string]interface{}{
			consulServiceUnregisterParam: definition.Id,
		},
		Body: nil,
	}
	client := consul.client

	request, err := http.NewRequest(requestConfig)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO UNREGISTER SERVICE: %s`, response.Status)
	}

	return nil
}

func watchServices(consul *Consul, index int) (int, error) {
	logger := consul.logger
	config := consul.config
	ctx := consul.ctx
	client := consul.client

	requestUrl := config.Url()
	requestUrl.Path = consulWatchServicesPath
	query := requestUrl.Query()
	query.Add(consulBlockingIndexKey, strconv.Itoa(index))
	query.Add(consulBlockingWaitKey, consulBlockingWaitValue)

	requestConfig := http.RequestConfig{
		Verb:    http.VerbGet,
		Url:     requestUrl,
		Headers: nil,
		Params:  nil,
		Body:    nil,
	}

	request, err := http.NewRequest(requestConfig)
	if err != nil {
		return -1, err
	}

	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return -1, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	headers := response.Header
	newIndex, err := strconv.ParseInt(headers.Get(consulBlockingHeaderKey), 10, 0)
	if err != nil {
		return -1, err
	}

	return int(newIndex), nil
}

func getServices(consul *Consul) (map[string][]string, error) {
	logger := consul.logger
	config := consul.config
	ctx := consul.ctx
	client := consul.client

	requestUrl := config.Url()
	requestUrl.Path = consulGetServicesPath
	requestConfig := http.RequestConfig{
		Verb:    http.VerbGet,
		Url:     requestUrl,
		Headers: nil,
		Params:  nil,
		Body:    nil,
	}

	request, err := http.NewRequest(requestConfig)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, errors.New(`NO SERVICES FOUND`)
	}

	services := make(map[string]interface{})
	err = json.Unmarshal(data, &services)
	if err != nil {
		return nil, err
	}

	updated := make(map[string][]string)
	for _, instance := range services {
		meta, ok := instance.(map[string]interface{})
		if !ok {
			continue
		}

		serviceRaw, ok := meta[`Service`]
		if !ok {
			continue
		}

		serviceName, ok := serviceRaw.(string)
		if !ok {
			continue
		}

		addressRaw, ok := meta[`Address`]
		if !ok {
			continue
		}

		address, ok := addressRaw.(string)
		if !ok {
			continue
		}

		addresses, ok := updated[serviceName]
		if !ok {
			updated[serviceName] = []string{address}
		} else {
			updated[serviceName] = append(addresses, address)
		}
	}

	return updated, nil
}

type Consul struct {
	logger   types.Logger
	config   Config
	client   *nativeHttp.Client
	ctx      context.Context
	cancel   func()
	services map[string][]string
	ready    chan bool
	mux      sync.Mutex
}

func (consul *Consul) Register(definition service.Service) error {
	logger := consul.logger
	err := registerService(consul, definition)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(`SERVICE %s REGISTERED SUCCESSFULLY`, definition.Name))
	return nil
}

func (consul *Consul) Unregister(definition service.Service) error {
	logger := consul.logger
	err := unregisterService(consul, definition)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(`SERVICE %s UNREGISTERED SUCCESSFULLY`, definition.Name))
	return nil
}

func (consul *Consul) Resolve(service string) ([]string, error) {
	ready := consul.ready

	var lastErr error = nil
	for i := 0; i < resolveRetry; i++ {
		_, ok := <-ready
		if ok {
			continue
		}

		services := consul.services
		addresses, ok := services[service]
		if !ok {
			lastErr = errors.Errorf(`SERVICE %s NOT FOUND`, service)
			if i+1 != resolveRetry {
				time.Sleep(resolveWait)
			}
			continue
		}

		if len(addresses) < 1 {
			lastErr = errors.Errorf(`SERVICE %s HAS NO ADDRESSES`, service)
			if i+1 != resolveRetry {
				time.Sleep(resolveWait)
			}
			continue
		}

		return addresses, nil
	}

	return []string{}, lastErr
}

func (consul *Consul) watchServices() {
	logger := consul.logger
	ctx := consul.ctx
	ready := consul.ready

	trackIndex := 0

	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil && err != context.Canceled {
				logger.Error(err)
				return
			}
		default:
			updatedIndex, err := watchServices(consul, trackIndex)
			if err != nil {
				if err == context.Canceled {
					return
				} else {
					logger.Error(err)
					time.Sleep(5 * time.Second)
					continue
				}
			}
			if updatedIndex == trackIndex {
				time.Sleep(5 * time.Second)
				continue
			}

			if trackIndex == 0 {
				close(ready)
			}
			trackIndex = updatedIndex

			services, err := getServices(consul)
			if err != nil {
				if err == context.Canceled {
					return
				} else {
					logger.Error(err)
					continue
				}
			}
			mux := consul.mux
			mux.Lock()
			consul.services = services
			mux.Unlock()

			logger.Info(fmt.Sprintf("AVAILABLE SERVICES: %+v", consul.services))
			time.Sleep(5 * time.Second)
		}
	}

}
