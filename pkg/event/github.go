package event

import (
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
		event := handler.config.FindMatchingEvent("Github", "comment", source)
		if event != nil {
			tasks = append(tasks, event.NewTask(source))
		}
	}
	return tasks
}
