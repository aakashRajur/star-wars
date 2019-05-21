package resource

import (
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/validate/normalizations"
	"github.com/aakashRajur/star-wars/pkg/validate/validations"
)

const (
	TypeResourceDefinition = `RESOURCE_DEFINITION`

	TypeString = `STRING`
	TypeInt    = `INT`

	resourceDefinitionHttpUri      = `HTTP_URI`
	resourceDefinitionHttpVerb     = `HTTP_VERB`
	resourceDefinitionType         = `TYPE`
	resourceDefinitionArgs         = `ARGS`
	resourceDefinitionArgKey       = `KEY`
	resourceDefinitionArgType      = `TYPE`
	resourceDefinitionArgRequired  = `REQUIRED`
	resourceDefinitionDataRequired = `DATA_REQUIRED`
)

type Arg struct {
	Key      string `json:"key" xml:"key"`
	Type     string `json:"type" xml:"type"`
	Required bool   `json:"required" xml:"required"`
}

type Definition struct {
	HttpURI      string `json:"http_uri" xml:"http_uri"`
	HttpVerb     string `json:"http_verb" xml:"http_verb"`
	Type         string `json:"type" xml:"type"`
	Args         []Arg  `json:"args" xml:"args"`
	DataRequired bool   `json:"data_required" xml:"data_required"`
	Source       string `json:"source" xml:"source"`
	AccessURI    string `json:"access_uri" xml:"access_uri"`
	Protocol     string `json:"protocol" xml:"protocol"`
	Endpoint     string `json:"endpoint" xml:"endpoint"`
}

func (resourceDefinition Definition) GetMap() map[string]interface{} {
	marshaled := make(map[string]interface{})

	marshaled[resourceDefinitionHttpUri] = resourceDefinition.HttpURI
	marshaled[resourceDefinitionHttpVerb] = resourceDefinition.HttpVerb
	marshaled[resourceDefinitionType] = resourceDefinition.Type

	args := make([]map[string]interface{}, 0)
	for _, arg := range resourceDefinition.Args {
		m := make(map[string]interface{})
		m[resourceDefinitionArgKey] = arg.Key
		m[resourceDefinitionArgType] = arg.Type
		m[resourceDefinitionArgRequired] = arg.Required

		args = append(args, m)
	}
	marshaled[resourceDefinitionArgs] = args

	marshaled[resourceDefinitionDataRequired] = resourceDefinition.DataRequired

	return marshaled
}

func (resourceDefinition Definition) GetArgValidators() map[string][]types.Validator {
	validators := make(map[string][]types.Validator)

	for _, each := range resourceDefinition.Args {
		compiled := make([]types.Validator, 0)
		if each.Required {
			compiled = append(compiled, validations.ValidateRequired())
		}
		switch each.Type {
		case TypeString:
			compiled = append(compiled, validations.ValidateString())
			break
		case TypeInt:
			compiled = append(compiled, validations.ValidateInteger())
		default:
			continue
		}

		validators[each.Key] = compiled
	}

	return validators
}

func (resourceDefinition Definition) GetArgNormalizers() map[string]types.Normalizor {
	normalizors := make(map[string]types.Normalizor)

	for _, each := range resourceDefinition.Args {
		switch each.Type {
		case TypeString:
			normalizors[each.Key] = normalizations.NormalizeString()
			break
		case TypeInt:
			normalizors[each.Key] = normalizations.NormalizeInteger()
		default:
			continue
		}
	}

	return normalizors
}

func (resourceDefinition Definition) Copy() Definition {
	return Definition{
		HttpURI:      resourceDefinition.HttpURI,
		HttpVerb:     resourceDefinition.HttpVerb,
		Type:         resourceDefinition.Type,
		Args:         resourceDefinition.Args,
		DataRequired: resourceDefinition.DataRequired,
		Source:       resourceDefinition.Source,
		AccessURI:    resourceDefinition.AccessURI,
		Protocol:     resourceDefinition.Protocol,
	}
}

func (Definition) Key() string {
	return TypeResourceDefinition
}
