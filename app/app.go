package app

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/tiny-job/core/shared"
)

type Option func(*options)

type options struct {
	keepAlive bool
	logger    hclog.Logger
}

// KeepAlive test use
func KeepAlive(keep bool) Option {
	return func(o *options) {
		o.keepAlive = keep
	}
}
func Logger(logger hclog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

type App struct {
	job  shared.Job
	opts options
}

func NewApp(job shared.Job, opts ...Option) *App {
	opt := options{
		keepAlive: false,
		logger: hclog.New(&hclog.LoggerOptions{
			Name:       "plugin",
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		}),
	}
	for _, o := range opts {
		o(&opt)
	}
	return &App{
		job:  job,
		opts: opt,
	}
}

func (a *App) Serve() {
	// 不需要父程序调用，就可以执行
	if a.opts.keepAlive {
		err := os.Setenv(shared.PluginJobEnv, shared.PluginJob)
		if err != nil {
			panic(err)
		}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginJob: &shared.JobRunPlugin{Impl: a.job},
		},
		Logger:     a.opts.logger,
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
