package consul

// Client ...
type Client interface {
	// Get a service from consul
	GetService(service, tag string) ([]string, error)
	// Register a service with local agent
	RegisterService(name, address string, tags []string, port int) error
	// Deregister a service with local agent
	DeRegisterService(id string) error
}
