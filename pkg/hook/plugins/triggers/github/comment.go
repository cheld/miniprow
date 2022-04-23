package github

import (
	"strings"

	config "github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/plugins/triggers"
	"github.com/go-playground/webhooks/v6/github"
)

const (
	HANDLER_ID = "github_comment"
)

func init() {
	triggers.RegisterHandler(HANDLER_ID, handleEvent)
}

func handleEvent(event config.Event, rule config.Rule) bool {
	searchText := rule.If.When["contains"]
	eventBody := event.Data.(github.IssueCommentPayload).Comment.Body
	ruleApplies := strings.Contains(strings.ToUpper(searchText), strings.ToUpper(eventBody))
	return ruleApplies
}
