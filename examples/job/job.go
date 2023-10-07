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
	log    hclog.Logger
}

func (j Job) Run(ctx context.Context, params map[string]string) (map[string]string, error) {
	j.log.Info("程序C运行开始 config:%+v", j.config)
	if params == nil {
		params = make(map[string]string)
	}
	params["C"] = "C"
	j.log.Info("程序C运行结束")
	return params, nil
}

var conf string

func init() {
	flag.StringVar(&conf, "conf", "", "config base64 values")
}

func main() {
	flag.Parse()

	config := getConfig()

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "plugin",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	// 启动服务
	app.NewApp(Job{config: config, log: logger}, app.Logger(logger)).Serve()
}

func getConfig() map[string]any {
	var config map[string]any
	if conf != "" {
		bytes, err := base64.StdEncoding.DecodeString(conf)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}
