package consul

func NewInstance(consulHost string) (*Consul, error) {
	instance := Consul{
		Host: consulHost,
	}

	return &instance, nil
}
