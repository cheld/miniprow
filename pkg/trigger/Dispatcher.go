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

func (dispatcher *Dispatcher) Execute(triggerInput []config.TriggerInput) {
	for _, triggerInput := range triggerInput {
		trigger := dispatcher.targets[triggerInput.Name]
		fmt.Println(trigger)
		switch trigger.Type {
		case "debug":
			ExecuteDebug(trigger, triggerInput)
		case "http":
			ExecuteHttp(trigger, triggerInput)
		default:
			fmt.Printf("I don't know about type %s!\n", "")
		}
	}

}
