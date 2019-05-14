package service

import (
	"encoding/json"
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

func (definition *Resource) WithProtocol(protocol string) *Resource {
	resource := Resource{
		ApiPattern: definition.ApiPattern,
		HttpVerb:   definition.HttpVerb,
		Type:       definition.Type,
		Protocol:   protocol,
	}
	return &resource
}
