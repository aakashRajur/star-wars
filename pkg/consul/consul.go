package consul

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

const (
	RegistrationUri   = `/v1/kv/:key`
	UnregistrationUri = `/v1/kv/:key`
)

func getResourceKey(resourceType string, source string) string {
	return fmt.Sprintf(`resources/%s/%s`, resourceType, source)
}

type Consul struct {
	Host string
}

func (consul *Consul) Register(definition resource_definition.ResourceDefinition) error {
	registerUrl := url.URL{
		Scheme: `http`,
		Host:   consul.Host,
		Path:   RegistrationUri,
	}
	registerConfig := http.RequestConfig{
		Verb: http.VerbPut,
		Url:  registerUrl,
		Params: map[string]interface{}{
			`key`: getResourceKey(definition.Type, definition.Source),
		},
		Body:    definition,
		Timeout: 10 * time.Second,
	}

	resp, err := http.NewRequest(registerConfig)
	if err != nil {
		return err
	}

	data, err := http.TextFromResponse(resp)
	if err != nil {
		return err
	}

	var success bool
	err = json.Unmarshal([]byte(data), &success)
	if err != nil {
		return err
	}

	if !success {
		return errors.Errorf(`UNABLE TO REGISTER RESOURCE %s`, definition.Type)
	}

	return nil
}

func (consul *Consul) Unregister(definition resource_definition.ResourceDefinition) error {
	unregisterUrl := url.URL{
		Scheme: `http`,
		Host:   consul.Host,
		Path:   UnregistrationUri,
	}
	unregisterConfig := http.RequestConfig{
		Verb: http.VerbPut,
		Url:  unregisterUrl,
		Params: map[string]interface{}{
			`key`: getResourceKey(definition.Type, definition.Source),
		},
		Timeout: 10 * time.Second,
	}

	resp, err := http.NewRequest(unregisterConfig)
	if err != nil {
		return err
	}

	data, err := http.TextFromResponse(resp)
	if err != nil {
		return err
	}

	var success bool
	err = json.Unmarshal([]byte(data), &success)
	if err != nil {
		return err
	}

	if !success {
		return errors.Errorf(`UNABLE TO UNREGISTER RESOURCE %s`, definition.Type)
	}

	return nil
}
