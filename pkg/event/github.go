package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
	"github.com/golang/glog"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	Github  = "Github"
	Comment = "Comment"
)

func (handler *Handler) HandleGithub(payload interface{}) []config.Task {
	var sourceName string
	var sourceType string

	// parse payload
	source := config.Source{
		Payload: payload,
		Environ: handler.env,
	}
	switch payload.(type) {
	case github.IssueCommentPayload:
		source.Value = payload.(github.IssueCommentPayload).Comment.Body
		sourceName = Github
		sourceType = Comment
	default:
		glog.Infof("Github event not implemented: %v\n", payload)
		return []config.Task{}
	}

	// handle event
	event := handler.config.FindEvent(sourceName, sourceType, source)
	if event == nil {
		glog.Infof("No event found for value %s", source.Value)
		return []config.Task{}
	}

	// build execution task
	task, err := event.BuildTask(source)
	if err != nil {
		glog.Errorf("Cannot handle event: %v. Error: %v", event.Trigger, err)
		return []config.Task{}
	}
	return []config.Task{task}
}
