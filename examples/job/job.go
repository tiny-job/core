package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/tiny-job/core/app"
)

type Job struct {
	config map[string]any
	logger hclog.Logger
}

func (j Job) Run(ctx context.Context, params map[string]string) (map[string]string, error) {
	j.logger.Info("程序C运行开始")
	if params == nil {
		params = make(map[string]string)
	}
	params["C"] = "C"
	j.logger.Info("程序C运行结束")
	return params, nil
}

var conf = flag.String("conf", "", "config base64 values")

func main() {
	flag.Parse()

	var config map[string]any
	if *conf != "" {
		bytes, err := base64.StdEncoding.DecodeString(*conf)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "plugin",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	// 启动服务
	app.NewApp(Job{logger: logger, config: config}, app.Logger(logger)).Serve()
}
