package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
)

func (handler *Handler) HandleCli(args, stdin string) []config.TriggerInput {
	triggerInput := []config.TriggerInput{}
	for _, event := range handler.config.Events {
		eventInput := config.EventInput{args, nil}
		if event.Type == "cli" && event.IsMatching(eventInput) {
			triggerInput = append(triggerInput, event.Handle(eventInput))
		}
	}
	return triggerInput
}
