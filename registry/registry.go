package registry

type Registry interface {
	Register(*Service) error
	Deregister(id string) error
	GetService(name string, tags ...string) ([]*Service, error)
}

type Service struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Tags     []string          `json:"tags"`
	Endpoint *Endpoint         `json:"endpoint"`
	Metadata map[string]string `json:"metadata"`
}

type Endpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
