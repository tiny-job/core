package balancer

import "github.com/tiny-job/core/registry"

type Balancer interface {
	Add(params ...*registry.Service) error
	Next() *registry.Service
}
