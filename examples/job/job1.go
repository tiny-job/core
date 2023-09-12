package main

import (
	"context"
	"log"

	"github.com/tiny-job/core/app"
	"github.com/tiny-job/core/registry/consul"
)

type Job01 struct {
}

func (Job01) Run(ctx context.Context, params map[string]string) (map[string]string, error) {
	log.Printf("欢迎使用")
	params["2"] = "2"
	return params, nil
}

func main() {
	// 使用consul服务作为发现
	registry := consul.NewRegistry("127.0.0.1:8500")
	// 启动服务
	app.NewApp(registry, Job01{}, app.KeepAlive(true)).Serve()
}
