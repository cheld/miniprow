package github

import (
	"strings"

	"github.com/cheld/miniprow/pkg/piper/config"
	"github.com/cheld/miniprow/pkg/piper/triggers"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	triggerName = "github_comment"
)

func init() {
	triggers.RegisterHandler(triggerName, handleEvent)
}

func handleEvent(event config.Event, rule config.Rule) bool {
	searchText := rule.If.When["contains"]
	eventBody := event.Data.(github.IssueCommentPayload).Comment.Body
	ruleApplies := strings.Contains(strings.ToUpper(searchText), strings.ToUpper(eventBody))
	return ruleApplies
}
