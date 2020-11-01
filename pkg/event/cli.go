package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
)

func (handler *Handler) HandleStdin(args, stdin string) []config.EventData {
	eventData := []config.EventData{}
	for _, rule := range handler.config.Events.Cli.Stdin {
		cliInput := config.InputData{args, nil}
		if rule.IsMatching(cliInput) {
			eventData = append(eventData, rule.Apply(cliInput))
		}
	}
	return eventData
}
