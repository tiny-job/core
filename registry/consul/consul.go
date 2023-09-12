package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/tiny-job/core/registry"
)

type consulRegistry struct {
	client       *api.Client
	config       *api.Config
	queryOptions *api.QueryOptions
}

func NewRegistry(addr string) registry.Registry {
	cr := &consulRegistry{
		queryOptions: &api.QueryOptions{
			AllowStale: true,
		},
	}

	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil
	}

	cr.client = client

	return cr
}

func (r consulRegistry) Register(s *registry.Service) error {
	// 1.注册的服务配置
	reg := api.AgentServiceRegistration{
		ID:      s.ID,
		Name:    s.Name,
		Tags:    s.Tags,
		Meta:    s.Metadata,
		Address: s.Endpoint.Host,
		Port:    s.Endpoint.Port,
		Check: &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", s.Endpoint.Host, s.Endpoint.Port),
			Timeout:                        "1s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}
	// 2.注册grpc服务到consul上
	err := r.client.Agent().ServiceRegister(&reg)
	if err != nil {
		return err
	}

	return nil
}

func (r consulRegistry) Deregister(id string) error {
	return r.client.Agent().ServiceDeregister(id)
}

func (r consulRegistry) GetService(name string, tags ...string) ([]*registry.Service, error) {
	entries, _, err := r.client.Health().ServiceMultipleTags(name, tags, true, r.queryOptions)
	if err != nil {
		return nil, err
	}

	var data []*registry.Service

	for _, entry := range entries {
		data = append(data, &registry.Service{
			ID:   entry.Service.ID,
			Name: entry.Service.Service,
			Tags: entry.Service.Tags,
			Endpoint: &registry.Endpoint{
				Host: entry.Service.Address,
				Port: entry.Service.Port,
			},
			Metadata: entry.Service.Meta,
		})
	}

	return data, err
}
