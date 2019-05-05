package resource_definition

import (
	"github.com/juju/errors"

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

type ResourceDefinition struct {
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

func (resourceDefinition ResourceDefinition) GetMap() map[string]interface{} {
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

func (resourceDefinition ResourceDefinition) GetArgValidators() map[string][]types.Validator {
	validators := make(map[string][]types.Validator)

	for _, each := range resourceDefinition.Args {
		compiled := make([]types.Validator, 0)
		if each.Required {
			compiled = append(compiled, validations.Required())
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

func (resourceDefinition ResourceDefinition) GetArgNormalizers() map[string]types.Normalizor {
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

func (resourceDefinition ResourceDefinition) Copy() ResourceDefinition {
	return ResourceDefinition{
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

func (ResourceDefinition) Key() string {
	return TypeResourceDefinition
}

func ResourceDefinitionFromMap(data map[string]interface{}) (ResourceDefinition, error) {
	definition := ResourceDefinition{}

	httpUri, ok := data[resourceDefinitionHttpUri]
	if !ok {
		return definition, errors.Errorf(`%s is required`, resourceDefinitionHttpUri)
	}
	definition.HttpURI, ok = httpUri.(string)
	if !ok {
		return definition, errors.Errorf(`%s should be a string`, resourceDefinitionHttpUri)
	}

	httpVerb, ok := data[resourceDefinitionHttpVerb]
	if !ok {
		return definition, errors.Errorf(`%s is required`, resourceDefinitionHttpVerb)
	}
	definition.HttpVerb, ok = httpVerb.(string)
	if !ok {
		return definition, errors.Errorf(`%s should be a string`, resourceDefinitionHttpVerb)
	}

	resourceType, ok := data[resourceDefinitionType]
	if !ok {
		return definition, errors.Errorf(`%s is required`, resourceDefinitionType)
	}
	definition.Type, ok = resourceType.(string)
	if !ok {
		return definition, errors.Errorf(`%s should be a string`, resourceDefinitionType)
	}

	jsonArgs, ok := data[resourceDefinitionArgs]
	if !ok {
		return definition, errors.Errorf(`%s is required`, resourceDefinitionArgs)
	}
	mapArgsArray, ok := jsonArgs.([]map[string]interface{})
	if !ok {
		return definition, errors.Errorf(
			`%s should be an array of {%s: STRING, %s: STRING, %s: BOOL}`,
			resourceDefinitionArgKey,
			resourceDefinitionArgType,
			resourceDefinitionArgRequired,
		)
	}

	parsedArgs := make([]Arg, 0)
	for i, each := range mapArgsArray {
		arg := Arg{}

		argKey, ok := each[resourceDefinitionArgKey]
		if !ok {
			return definition, errors.Errorf(`%s at index %d is required`, resourceDefinitionArgKey, i)
		}
		arg.Key, ok = argKey.(string)
		if !ok {
			return definition, errors.Errorf(`%s at index %d should be a string`, resourceDefinitionArgKey, i)
		}

		argType, ok := each[resourceDefinitionArgType]
		if !ok {
			return definition, errors.Errorf(`%s at index %d is required`, resourceDefinitionArgType, i)
		}
		arg.Type, ok = argType.(string)
		if !ok {
			return definition, errors.Errorf(`%s at index %d should be a string`, resourceDefinitionArgType, i)
		}

		argRequired, ok := each[resourceDefinitionArgRequired]
		if !ok {
			return definition, errors.Errorf(`%s at index %d is required`, resourceDefinitionArgRequired, i)
		}
		arg.Required, ok = argRequired.(bool)
		if !ok {
			return definition, errors.Errorf(`%s at index %d should be a bool`, resourceDefinitionArgRequired, i)
		}

		parsedArgs = append(parsedArgs, arg)
	}

	dataRequired, ok := data[resourceDefinitionDataRequired]
	if !ok {
		return definition, errors.Errorf(`%s is required`, resourceDefinitionDataRequired)
	}
	definition.DataRequired, ok = dataRequired.(bool)
	if !ok {
		return definition, errors.Errorf(`%s should be a bool`, resourceDefinitionType)
	}

	return definition, nil
}
