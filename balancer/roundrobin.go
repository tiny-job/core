package balancer

import (
	"errors"

	"github.com/tiny-job/core/registry"
)

// RoundRobinBalance 轮询负载均衡
type RoundRobinBalance struct {
	curIndex int
	rss      []*registry.Service
}

func (r *RoundRobinBalance) Add(params ...*registry.Service) error {
	if len(params) == 0 {
		return errors.New("params len 1 at least")
	}

	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RoundRobinBalance) Next() *registry.Service {
	if len(r.rss) == 0 {
		return nil
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}
