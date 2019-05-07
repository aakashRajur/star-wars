package resource

type Registree interface {
	Register(definition Definition) error
	Unregister(definition Definition) error
}
