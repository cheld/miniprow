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

func (handler *Handler) HandleGithub(payload interface{}) config.Ctx {

	// parse payload
	request := config.RequestCtx{
		Payload: payload,
	}
	ctx := config.Ctx{
		Request: request,
		Environ: *util.Environment.Map(),
	}
	switch payload.(type) {
	case github.IssueCommentPayload:
		ctx.Request.Value = payload.(github.IssueCommentPayload).Comment.Body
		ctx.Request.Event = "github_comment"
	default:
		glog.Infof("Github event not implemented: %v\n", payload)
		return ctx
	}

	// handle rule
	handler.config.ApplyRule(&ctx)

	return ctx
}
