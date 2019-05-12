package consul

import (
	"fmt"
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
	consulResourceKey            = `key`
	consulResourceRegisterPath   = `/v1/kv/:key`
	consulResourceUnregisterPath = `/v1/kv/:key`
	resolveRetry                 = 5
	resolveWait                  = 5 * time.Second
)

func registerService(consul *Consul, definition service.Service) error {
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

	url := consul.Config.Url()
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
	if response.StatusCode < 200 && response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO REGISTER SERVICE: %s`, response.Status)
	}

	return nil
}

func unregisterService(consul *Consul, definition service.Service) error {
	url := consul.Config.Url()
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

	if response.StatusCode < 200 && response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO UNREGISTER SERVICE: %s`, response.Status)
	}

	return nil
}

func registerResources(consul *Consul, definition service.Service) {
	logger := consul.Logger
	resources := definition.Resources
	for _, resource := range resources {
		url := consul.Config.Url()
		url.Path = consulResourceRegisterPath

		body := resource.Map()

		requestConfig := http.RequestConfig{
			Verb:    http.VerbPut,
			Url:     url,
			Headers: nil,
			Params: map[string]interface{}{
				consulResourceKey: fmt.Sprintf(
					`%s/%s`,
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
		if response.StatusCode < 200 && response.StatusCode > 299 {
			logger.Error(errors.Errorf(`UNABLE TO REGISTER RESOURCE %s`, resource.Type))
		} else {
			logger.Info(fmt.Sprintf(`RESOURCE %s SUCCESSFULLY`, resource.Type))
		}
	}
}

func unregisterResources(consul *Consul, definition service.Service) {
	logger := consul.Logger
	resources := definition.Resources
	for _, resource := range resources {
		resource.Protocol = definition.Scheme
		url := consul.Config.Url()
		url.Path = consulResourceUnregisterPath

		requestConfig := http.RequestConfig{
			Verb:    http.VerbDelete,
			Url:     url,
			Headers: nil,
			Params: map[string]interface{}{
				consulResourceKey: fmt.Sprintf(
					`%s/%s`,
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
		if response.StatusCode < 200 && response.StatusCode > 299 {
			logger.Error(errors.Errorf(`UNABLE TO UNREGISTER RESOURCE %s`, resource.Type))
		} else {
			logger.Info(fmt.Sprintf(`RESOURCE %s UNREGISTERED SUCCESSFULLY`, resource.Type))
		}
	}
}

type Consul struct {
	Logger types.Logger
	Config Config
}

func (consul *Consul) Register(definition service.Service) error {
	logger := consul.Logger
	err := registerService(consul, definition)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(`SERVICE %s REGISTERED SUCCESSFULLY`, definition.Name))
	registerResources(consul, definition)
	return nil
}

func (consul *Consul) Unregister(definition service.Service) error {
	logger := consul.Logger
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
		url := consul.Config.Url()
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
