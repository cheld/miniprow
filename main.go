package main

import (
	"github.com/cheld/cicd-bot/pkg/config"
	"github.com/cheld/cicd-bot/pkg/event"
	"github.com/cheld/cicd-bot/pkg/trigger"
	"github.com/cheld/cicd-bot/pkg/webhook"
)

func main() {

	config := config.Load("config.yaml")
	args := "test"
	stdin := ""
	triggerInput := event.NewHandler(config).HandleCli(args, stdin)
	trigger.NewDispatcher(config).Execute(triggerInput)

	webhook.Run()
}
