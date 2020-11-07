package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
)

func (handler *Handler) HandleCli(args, stdin string) []config.Task {
	tasks := []config.Task{}
	for _, event := range handler.config.Events {
		source := config.Source{args, nil}
		if event.Source == "cli" && event.IsMatching(source) {
			tasks = append(tasks, event.NewTask(source))
		}
	}
	return tasks
}
