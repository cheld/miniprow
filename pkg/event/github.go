package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func (handler *Handler) HandleGithub(payload interface{}) []config.TriggerInput {

	eventInput := config.EventInput{
		Objectiv: "",
		Input: map[string]interface{}{ // TODO remove. Name payload
			"Payload": payload,
		},
	}

	triggerInput := []config.TriggerInput{}

	switch payload.(type) {
	case github.IssueCommentPayload:
		commentPayload := payload.(github.IssueCommentPayload)
		eventInput.Objectiv = commentPayload.Comment.Body
		event := handler.config.FindMatchingEvent("Github", "comment", eventInput)
		if event != nil {
			triggerInput = append(triggerInput, event.Handle(eventInput))
		}
	}
	return triggerInput
}
