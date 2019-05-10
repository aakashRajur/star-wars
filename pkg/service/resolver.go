package service

type Resolver interface {
	Resolve(service string) ([]string, error)
}
