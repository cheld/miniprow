package github

import (
	"github.com/cheld/miniprow/pkg/piper/action/actions"
	"github.com/cheld/miniprow/pkg/piper/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	actionName = "github_comment"
)

func init() {
	actions.RegisterHandler(actionName, handleAction)
}

func handleAction(rule config.Rule, event config.Event) {
	issueNumber := event.Data.(github.IssueCommentPayload).Issue.Number
}
