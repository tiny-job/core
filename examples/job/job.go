package main

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/tiny-job/core/app"
)

type Job struct {
	logger hclog.Logger
}

func (j Job) Run(ctx context.Context, params map[string]string) (map[string]string, error) {
	j.logger.Debug("欢迎使用程序B")
	params["B"] = "B"
	return params, nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "plugin",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	// 启动服务
	app.NewApp(Job{logger: logger}, app.Logger(logger)).Serve()
}
