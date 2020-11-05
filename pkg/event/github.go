package event

import (
	"fmt"

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
		fmt.Println(eventInput.Objectiv)
		for _, event := range handler.config.Events {
			if event.Source == "github" && event.Type == "comment" && event.IsMatching(eventInput) {
				triggerInput = append(triggerInput, event.Handle(eventInput))
			}
		}
	}
	return triggerInput
}
