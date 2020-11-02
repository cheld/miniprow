package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
)

func (handler *Handler) HandleCli(args, stdin string) []config.EventData {
	eventData := []config.EventData{}
	for _, eventSpec := range handler.config.Events.Cli {
		cliInput := config.InputData{args, nil}
		if eventSpec.IsMatching(cliInput) {
			eventData = append(eventData, eventSpec.Process(cliInput))
		}
	}
	return eventData
}
