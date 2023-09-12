package client

import (
	"context"
	"errors"
	"net"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cast"
	"github.com/tiny-job/core/balancer"
	"github.com/tiny-job/core/registry"
	"github.com/tiny-job/core/selector"
	"github.com/tiny-job/core/shared"
)

type Option func(*options)

type options struct {
	name string   // 任务名称
	tags []string // 任务标签
}

var defaultOptions = options{
	tags: []string{shared.TagConsul, shared.TagGRPC, shared.TagJob},
}

func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func Tags(tags ...string) Option {
	return func(o *options) {
		o.tags = tags
	}
}

type Client struct {
	opts options

	register registry.Registry

	logger hclog.Logger
}

func NewClient(register registry.Registry, logger hclog.Logger, opts ...Option) *Client {
	opt := defaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return &Client{
		opts:     opt,
		register: register,
		logger:   logger,
	}
}

type Runner struct {
	client *plugin.Client
}

func (c *Client) Plugin(random bool) (*Runner, error) {

	services, err := c.register.GetService(c.opts.name, c.opts.tags...)
	if err != nil {
		return nil, err
	}

	var bcr balancer.Balancer
	if random {
		b, err := selector.Random(services)
		if err != nil {
			return nil, err
		}
		bcr = b
	} else {
		b, err := selector.RoundRobin(services)
		if err != nil {
			return nil, err
		}
		bcr = b
	}

	srv := bcr.Next()
	if srv == nil {
		return nil, errors.New("service not found")
	}

	addrStr := net.JoinHostPort(srv.Endpoint.Host, strconv.Itoa(srv.Endpoint.Port))
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}

	pid := cast.ToInt(srv.Metadata[shared.MetaDataPid])
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Reattach: &plugin.ReattachConfig{
			Protocol: plugin.ProtocolGRPC,
			Addr:     addr,
			Pid:      pid,
			Test:     false,
		},
		// Cmd: exec.Command("./jobs/job1.exe"),
		Logger: c.logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})
	return &Runner{
		client: client,
	}, nil
}

func (r *Runner) Run(ctx context.Context, params map[string]string) (map[string]string, error) {
	rpcClient, err := r.client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.PluginJob)
	if err != nil {
		return nil, err
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	job := raw.(shared.Job)

	result, err := job.Run(ctx, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Runner) Kill() {
	r.client.Kill()
}
