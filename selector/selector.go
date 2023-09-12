package selector

import (
	"github.com/tiny-job/core/balancer"
	"github.com/tiny-job/core/registry"
)

func Random(s []*registry.Service) (balancer.Balancer, error) {
	random := balancer.RandomBalance{}
	err := random.Add(s...)
	if err != nil {
		return nil, err
	}
	return &random, nil
}

func RoundRobin(s []*registry.Service) (balancer.Balancer, error) {
	roundRobin := balancer.RoundRobinBalance{}
	err := roundRobin.Add(s...)
	if err != nil {
		return nil, err
	}
	return &roundRobin, nil
}
