package resource_definition

type ResourceRegistration interface {
	Register(definition ResourceDefinition) error
	Unregister(definition ResourceDefinition) error
}
