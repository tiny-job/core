package main

import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/tiny-job/core/client"
	"github.com/tiny-job/core/registry/consul"
	"golang.org/x/net/context"
)

func main() {
	// 使用consul服务作为发现
	registry := consul.NewRegistry("localhost:8500")
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "main",
		Output:     os.Stdout,
		JSONFormat: true,
	})
	// 启动服务
	cli := client.NewClient(registry, logger,
		client.Name("Job01"),
		client.Tags("consul", "grpc", "job"),
	)

	runner, err := cli.Plugin(true)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := runner.Run(context.Background(), map[string]string{"test": "test01"})
	if err != nil {
		log.Fatalln(err)
	}

	runner.Kill()

	log.Println(result)
}
