package event

import (
	"github.com/cheld/cicd-bot/pkg/piper/config"
	"github.com/golang/glog"
)

func (handler *Handler) HandleCli(eventType, payload string) []config.Task {
	glog.V(7).Infof("Cli event handler '%s' received payload:\n%s\n\n", eventType, payload)

	// parse input
	source := config.Source{payload, payload, handler.env}

	// handle event
	event := handler.config.FindEvent("cli", eventType, source)
	if event == nil {
		glog.Infof("No event handler found for cli, type %s\n", eventType)
		return []config.Task{}
	}

	// build execution task
	task, err := event.BuildTask(source)
	if err != nil {
		glog.Errorf("Cannot handle event: %v. Error: %v", event.Trigger, err)
		return []config.Task{}
	}
	glog.V(7).Infof("Execution task '%s' created. Values:\n%s\n\n", task.Trigger, task.Values)
	return []config.Task{task}
}
