package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/tiny-job/core/client"
	"golang.org/x/net/context"
)

func main() {
	// 使用consul服务作为发现
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "executor",
		Output:     os.Stdout,
		JSONFormat: true,
	})
	// 启动服务
	cli := client.NewClient(
		client.Logger(logger),
	)

	runner, err := cli.Plugin(exec.Command("./jobs/job.exe"))
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
