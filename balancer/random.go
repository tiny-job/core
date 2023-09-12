package balancer

import (
	"errors"
	"math/rand"

	"github.com/tiny-job/core/registry"
)

// RandomBalance 随机负载均衡
type RandomBalance struct {
	curIndex int

	rss []*registry.Service
}

func (r *RandomBalance) Add(params ...*registry.Service) error {
	if len(params) == 0 {
		return errors.New("params len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)

	return nil
}

func (r *RandomBalance) Next() *registry.Service {
	if len(r.rss) == 0 {
		return nil
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}
