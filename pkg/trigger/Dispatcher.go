package trigger

import (
	"strings"

	"github.com/cheld/cicd-bot/pkg/config"
	"github.com/golang/glog"
)

type Dispatcher struct {
	config config.Configuration
}

func NewDispatcher(cfg config.Configuration) Dispatcher {
	dispatcher := Dispatcher{}
	dispatcher.config = cfg
	return dispatcher
}

func (dispatcher *Dispatcher) Execute(tasks []config.Task) {
	for _, task := range tasks {
		trigger := dispatcher.config.FindTrigger(task.Trigger)
		if trigger == nil {
			glog.Errorf("No trigger definition with name '%s' found\n", task.Trigger)
			break
		}
		switch strings.ToLower(trigger.Type) {
		case "debug":
			ExecuteDebug(trigger, task)
		case "http":
			ExecuteHttp(trigger, task)
		default:
			glog.Errorf("No implementation for trigger type '%s' found!\n", trigger.Type)
		}
	}
}
