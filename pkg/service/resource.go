package service

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	keyApiEndpoint = `api_pattern`
	keyHttpVerb    = `http_verb`
	keyType        = `type`
	keyProtocol    = `protocol`
)

type Resource struct {
	ApiPattern string `json:"api_pattern" xml:"api_pattern"`
	HttpVerb   string `json:"http_verb" xml:"http_verb"`
	Type       string `json:"type" xml:"type"`
	Protocol   string `json:"protocol" xml:"protocol"`
}

func (definition *Resource) Map() map[string]interface{} {
	mapped := make(map[string]interface{})
	mapped[keyApiEndpoint] = definition.ApiPattern
	mapped[keyHttpVerb] = definition.HttpVerb
	mapped[keyType] = definition.Type
	mapped[keyProtocol] = definition.Protocol
	return mapped
}

func (definition *Resource) JSON() (string, error) {
	marshaled, err := json.Marshal(definition)
	if err != nil {
		return ``, err
	}
	return string(marshaled), nil
}

func ResourceDefinitionFromMap(mapped map[string]interface{}) (Resource, error) {
	rd := Resource{}

	apiEndpoint, ok := mapped[keyApiEndpoint]
	if !ok {
		return rd, errors.Errorf(`%s NOT PROVIDED`, keyApiEndpoint)
	}
	rd.ApiPattern, ok = apiEndpoint.(string)
	if !ok {
		return rd, errors.Errorf(`%s SHOULD BE A STRING`, keyApiEndpoint)
	}

	httpVerb, ok := mapped[keyHttpVerb]
	if !ok {
		return rd, errors.Errorf(`%s NOT PROVIDED`, keyHttpVerb)
	}
	rd.ApiPattern, ok = httpVerb.(string)
	if !ok {
		return rd, errors.Errorf(`%s SHOULD BE A STRING`, keyHttpVerb)
	}

	_type, ok := mapped[keyType]
	if !ok {
		return rd, errors.Errorf(`%s NOT PROVIDED`, keyType)
	}
	rd.ApiPattern, ok = _type.(string)
	if !ok {
		return rd, errors.Errorf(`%s SHOULD BE A STRING`, keyType)
	}

	return rd, nil
}
