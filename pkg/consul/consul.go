package consul

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	consulServiceQueryParam      = `service`
	consulServiceQueryPath       = `/v1/catalog/service/:service`
	consulServiceAddressKey      = `ServiceAddress`
	consulKeyPrefix              = `resources`
	consulResourceKey            = `key`
	consulResourceRegisterPath   = `/v1/kv/:key`
	consulResourceUnregisterPath = `/v1/kv/:key`
	consulResourceReadPath       = `/v1/kv/:key`
	consulIndexKey               = `X-Consul-Index`
	consulResourceReadIndexKey   = `index`
	consulResourcePollKey        = `wait`
	consulResourcePollValue      = `5m`
	resolveRetry                 = 5
	resolveWait                  = 5 * time.Second
)

func registerService(consul *Consul, definition service.Service) error {
	config := consul.config
	body := make(map[string]interface{})

	body[keyId] = definition.Id
	body[keyName] = definition.Name
	if definition.Port > 0 {
		body[keyAddress] = fmt.Sprintf(
			`%s://%s:%d`,
			definition.Scheme,
			definition.Hostname,
			definition.Port,
		)
	}
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
		Timeout: 10 * time.Second,
	}

	response, err := http.NewRequest(requestConfig)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO REGISTER SERVICE: %s`, response.Status)
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
		Body:    nil,
		Timeout: 10 * time.Second,
	}

	response, err := http.NewRequest(requestConfig)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO UNREGISTER SERVICE: %s`, response.Status)
	}

	return nil
}

func registerResources(consul *Consul, definition service.Service) {
	logger := consul.logger
	config := consul.config
	resources := definition.Resources

	for _, resource := range resources {
		url := config.Url()
		url.Path = consulResourceRegisterPath

		body := resource.Map()

		requestConfig := http.RequestConfig{
			Verb:    http.VerbPut,
			Url:     url,
			Headers: nil,
			Params: map[string]interface{}{
				consulResourceKey: fmt.Sprintf(
					`%s/%s/%s`,
					consulKeyPrefix,
					resource.Type,
					definition.Id,
				),
			},
			Body:    body,
			Timeout: 10 * time.Second,
		}

		response, err := http.NewRequest(requestConfig)
		if err != nil {
			logger.Error(err)
		}
		if response.StatusCode < 200 || response.StatusCode > 299 {
			logger.Error(errors.Errorf(`UNABLE TO REGISTER RESOURCE %s`, resource.Type))
		} else {
			logger.Info(fmt.Sprintf(`RESOURCE %s SUCCESSFULLY`, resource.Type))
		}
	}
}

func unregisterResources(consul *Consul, definition service.Service) {
	logger := consul.logger
	config := consul.config
	resources := definition.Resources

	for _, resource := range resources {
		resource.Protocol = definition.Scheme
		url := config.Url()
		url.Path = consulResourceUnregisterPath

		requestConfig := http.RequestConfig{
			Verb:    http.VerbDelete,
			Url:     url,
			Headers: nil,
			Params: map[string]interface{}{
				consulResourceKey: fmt.Sprintf(
					`%s/%s/%s`,
					consulKeyPrefix,
					resource.Type,
					definition.Id,
				),
			},
			Body:    nil,
			Timeout: 10 * time.Second,
		}

		response, err := http.NewRequest(requestConfig)
		if err != nil {
			logger.Error(err)
		}
		statusCode := response.StatusCode
		if (statusCode < 200 || statusCode > 299) && statusCode != nativeHttp.StatusNotFound {
			logger.Error(errors.Errorf(`UNABLE TO UNREGISTER RESOURCE %s`, resource.Type))
		} else {
			logger.Info(fmt.Sprintf(`RESOURCE %s UNREGISTERED SUCCESSFULLY`, resource.Type))
		}
	}
}

func emit(consul *Consul, key string, data interface{}) {
	logger := consul.logger
	ctx := consul.ctx

	consul.mux.Lock()
	observations, ok := consul.observations[key]
	consul.mux.Unlock()
	if !ok {
		logger.Warn(fmt.Sprintf(`NO SUBSCRIPTIONS FOR %s`, key))
	}

	for _, each := range observations {
		go each.Handler(data, ctx)
	}
}

func watchKV(consul *Consul, subscription service.Subscription) {
	logger := consul.logger
	config := consul.config
	ctx := consul.ctx

	trackIndex := 0

	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				logger.Error(err)
				return
			}
		default:
			indexString := strconv.Itoa(trackIndex)
			url := config.Url()
			url.Path = consulResourceRegisterPath
			query := url.Query()
			query.Add(consulResourceReadIndexKey, indexString)
			query.Add(consulResourcePollKey, consulResourcePollValue)
			query.Add(`recurse`, `true`)

			requestConfig := http.RequestConfig{
				Verb: http.VerbGet,
				Url:  url,
				Headers: map[string]string{
					consulIndexKey: indexString,
				},
				Params: map[string]interface{}{
					consulResourceReadPath: subscription.Key,
				},
				Body: nil,
			}

			response, err := http.NewCanceleableRequest(requestConfig, ctx)
			if err != nil {
				if err == context.Canceled {
					return
				} else {
					logger.Error(err)
				}
			}

			statusCode := response.StatusCode
			if statusCode == nativeHttp.StatusNotFound {
				emit(consul, subscription.Key, nil)
				time.Sleep(5 * time.Second)
			} else if statusCode < 200 || statusCode > 299 {
				logger.Error(http.TextFromResponse(response))
				time.Sleep(5 * time.Second)
				return
			}

			headers := response.Header
			index, err := strconv.ParseInt(headers.Get(consulIndexKey), 10, 0)
			if err != nil {
				logger.Error(err)
			} else {
				trackIndex = int(index)
			}

			data, err := http.TextFromResponse(response)
			if err != nil {
				logger.Error(err)
				continue
			}

			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				logger.Error(err)
				continue
			}

			marshaled := make([]map[string]interface{}, 0)
			err = json.Unmarshal(decoded, &marshaled)
			if err != nil {
				logger.Error(err)
				continue
			}

			emit(consul, subscription.Key, marshaled)
		}
	}
}

type Consul struct {
	logger       types.Logger
	config       Config
	ctx          context.Context
	cancel       func()
	observations map[string][]service.Subscription
	mux          sync.Mutex
}

func (consul *Consul) Register(definition service.Service) error {
	logger := consul.logger
	err := registerService(consul, definition)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(`SERVICE %s REGISTERED SUCCESSFULLY`, definition.Name))
	registerResources(consul, definition)
	return nil
}

func (consul *Consul) Unregister(definition service.Service) error {
	logger := consul.logger
	unregisterResources(consul, definition)
	err := unregisterService(consul, definition)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(`SERVICE %s UNREGISTERED SUCCESSFULLY`, definition.Name))
	return nil
}

func (consul *Consul) Resolve(service string) ([]string, error) {
	var lastErr error = nil
	for i := 0; i < resolveRetry; i++ {
		url := consul.config.Url()
		url.Path = consulServiceQueryPath

		requestConfig := http.RequestConfig{
			Verb:    http.VerbGet,
			Url:     url,
			Headers: nil,
			Params: map[string]interface{}{
				consulServiceQueryParam: service,
			},
			Body:    nil,
			Timeout: 10 * time.Second,
		}

		compiled := make([]string, 0)
		response, err := http.NewRequest(requestConfig)
		if err != nil {
			return compiled, err
		}

		available, err := http.JsonArrayFromResponse(response)
		if err != nil {
			lastErr = err
			time.Sleep(resolveWait)
			continue
		}

		for _, each := range available {
			address, ok := (each[consulServiceAddressKey]).(string)
			if !ok {
				continue
			}
			compiled = append(compiled, address)
		}

		if len(compiled) < 1 {
			lastErr = errors.Errorf(`NO SERVICES FOUND FOR %s`, service)
			time.Sleep(resolveWait)
			continue
		}

		return compiled, nil
	}

	return []string{}, lastErr
}

func (consul *Consul) Observe(subscription service.Subscription) error {
	mux := consul.mux
	mux.Lock()
	defer mux.Unlock()

	observations := consul.observations
	subscriptions, ok := observations[subscription.Key]
	if ok {
		subscriptions = append(subscriptions, subscription)
	} else {
		subscriptions = []service.Subscription{subscription}
		go watchKV(consul, subscription)
	}
	observations[subscription.Key] = subscriptions

	return nil
}
