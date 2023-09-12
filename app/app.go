package app

import (
	"net"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cast"
	"github.com/tiny-job/core/registry"
	"github.com/tiny-job/core/shared"
)

type Option func(*options)

type options struct {
	id        string
	pid       string
	name      string
	host      string
	port      int
	keepAlive bool
}

var defaultOptions = options{
	id:        uuid.New().String(),
	keepAlive: false,
}

func ID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func KeepAlive(keep bool) Option {
	return func(o *options) {
		o.keepAlive = keep
	}
}

type App struct {
	register registry.Registry
	job      shared.Job
	opts     options

	logger hclog.Logger
}

func NewApp(register registry.Registry, job shared.Job, opts ...Option) *App {
	opt := defaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return &App{
		register: register,
		job:      job,
		opts:     opt,
	}
}

func (a *App) init() {
	a.opts.pid = cast.ToString(os.Getpid())
	if len(a.opts.id) == 0 {
		a.opts.id = uuid.New().String()
	}
	if len(a.opts.name) == 0 {
		jobType := reflect.TypeOf(a.job)
		a.opts.name = jobType.Name()
	}

	handler := func(level hclog.Level, msg string, args ...interface{}) bool {
		if level == hclog.Debug && msg == "plugin address" {
			// 	注册服务
			addr := args[3].(string)
			_, port, err := net.SplitHostPort(addr)
			if err != nil {
				panic(err)
			}

			a.opts.port = cast.ToInt(port)
			a.opts.host = "192.168.3.13"

			err = a.registerService()
			if err != nil {
				panic(err)
			}
			return false
		}
		return false
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "plugin",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
		Exclude:    handler,
	})
	a.logger = logger
}

func (a *App) Serve() {
	// 永久运行任务
	if a.opts.keepAlive {
		err := os.Setenv(shared.PluginJobEnv, shared.PluginJob)
		if err != nil {
			panic(err)
		}
	}

	a.init()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginJob: &shared.JobRunPlugin{Impl: a.job},
		},
		Logger:     a.logger,
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

func (a *App) registerService() error {
	return a.register.Register(&registry.Service{
		ID:   a.opts.id,
		Name: a.opts.name,
		Tags: []string{shared.TagConsul, shared.TagGRPC, shared.TagJob},
		Endpoint: &registry.Endpoint{
			Host: a.opts.host,
			Port: a.opts.port,
		},
		Metadata: map[string]string{
			shared.MetaDataPid: a.opts.pid,
		},
	})

}
