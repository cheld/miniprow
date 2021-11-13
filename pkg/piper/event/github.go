package event

import (
	"github.com/cheld/miniprow/pkg/common/util"
	"github.com/cheld/miniprow/pkg/piper/config"
	"github.com/golang/glog"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	Github  = "Github"
	Comment = "Comment"
)

func (handler *Handler) HandleGithub(payload interface{}) []config.Task {
	var eventName string

	// parse payload
	request := config.Request{
		Payload: payload,
	}
	ctx := config.Ctx{
		Request: request,
		Environ: *util.Environment.Map(),
	}
	switch payload.(type) {
	case github.IssueCommentPayload:
		ctx.Request.Value = payload.(github.IssueCommentPayload).Comment.Body
		eventName = "github_comment"
	default:
		glog.Infof("Github event not implemented: %v\n", payload)
		return []config.Task{}
	}

	// handle event
	event := handler.config.GetMatchingRule(eventName, ctx)
	if event == nil {
		glog.Infof("No event found for value %s", ctx.Request.Value)
		return []config.Task{}
	}

	// build execution task
	task, err := event.BuildTask(ctx)
	if err != nil {
		glog.Errorf("Cannot handle event: %v. Error: %v", event.Trigger, err)
		return []config.Task{}
	}
	return []config.Task{task}
}
