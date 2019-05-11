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
)

type Consul struct {
	Logger types.Logger
	Config Config
}

func (consul *Consul) Register(definition service.Service) error {
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

func (consul *Consul) Unregister(definition service.Service) error {
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
	fmt.Printf("UNREGISTER:  %+v\n", response)

	if response.StatusCode < 200 && response.StatusCode > 299 {
		return errors.Errorf(`FAILED TO UNREGISTER SERVICE: %s`, response.Status)
	}

	return nil
}

func (consul *Consul) Resolve(service string) ([]string, error) {
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
		return compiled, err
	}

	for _, each := range available {
		address, ok := (each[consulServiceAddressKey]).(string)
		if !ok {
			continue
		}
		compiled = append(compiled, address)
	}

	return compiled, nil
}
