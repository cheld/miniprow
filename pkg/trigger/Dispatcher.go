package trigger

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

type Dispatcher struct {
	config  config.Configuration
	targets map[string]config.Trigger
}

func NewDispatcher(cfg config.Configuration) *Dispatcher {
	dispatcher := Dispatcher{}
	dispatcher.config = cfg
	dispatcher.targets = make(map[string]config.Trigger)
	for _, trigger := range cfg.Triggers {
		dispatcher.targets[trigger.Name] = trigger
	}
	return &dispatcher
}

func (dispatcher *Dispatcher) Execute(tasks []config.Task) {
	for _, task := range tasks {
		trigger := dispatcher.targets[task.Trigger]
		switch trigger.Type {
		case "debug":
			ExecuteDebug(trigger, task)
		case "http":
			ExecuteHttp(trigger, task)
		default:
			fmt.Printf("I don't know about type %s!\n", "")
		}
	}

}
