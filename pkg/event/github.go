package event

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func (handler *Handler) HandleGithub(payload interface{}) []config.Task {
	source := config.Source{
		Payload: payload,
	}

	tasks := []config.Task{}

	switch payload.(type) {
	case github.IssueCommentPayload:
		commentPayload := payload.(github.IssueCommentPayload)
		source.Value = commentPayload.Comment.Body
		event := handler.config.Event("Github", "comment", source)
		if event == nil {
			fmt.Printf("No event found for value %s", source.Value)
		} else {
			task, err := event.NewTask(source)
			if err != nil {
				fmt.Printf("Cannot handle event: %v. Error: %v", event.Trigger, err)
			} else {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}
