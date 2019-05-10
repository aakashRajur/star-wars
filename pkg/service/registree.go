package service

type Registree interface {
	Register(definition Service) error
	Unregister(definition Service) error
}
