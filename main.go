package main

import (
	"github.com/cheld/cicd-bot/pkg/webhook"
)

func main() {

	//config := config.Load("config.yaml")
	//args := "test"
	//stdin := ""
	//tasks := event.NewHandler(config).HandleCli(args, stdin)
	//trigger.NewDispatcher(config).Execute(tasks)

	webhook.Run()
}
