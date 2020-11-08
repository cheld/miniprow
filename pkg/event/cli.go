package event

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

func (handler *Handler) HandleCli(args, stdin string) []config.Task {
	tasks := []config.Task{}
	for _, event := range handler.config.Events {
		source := config.Source{args, nil}
		if event.Source == "cli" && event.IsMatching(source) {
			task, err := event.NewTask(source)
			if err != nil {
				fmt.Errorf("Cannot trigger: %v. Error: %v", task.Trigger, err)
			} else {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}
