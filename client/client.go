package client

import (
	"context"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/tiny-job/core/shared"
)

type Option func(*options)

type options struct {
	logger hclog.Logger
}

var defaultOptions = options{
	logger: hclog.New(&hclog.LoggerOptions{
		Name:       "executor",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	}),
}

func Logger(logger hclog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

type Client struct {
	opts options
}

func NewClient(opts ...Option) *Client {
	opt := defaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return &Client{
		opts: opt,
	}
}

type Runner struct {
	client *plugin.Client
}

func (c *Client) Plugin(cmd *exec.Cmd) (*Runner, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             cmd,
		Logger:          c.opts.logger,
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
