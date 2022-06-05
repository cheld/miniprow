package github

import (
	"strings"

	"github.com/cheld/miniprow/pkg/common/core"
	"github.com/cheld/miniprow/pkg/hook/rules/filters"
	"github.com/go-playground/webhooks/v6/github"
)

const (
	HANDLER_ID = "github_comment"
)

func init() {
	filters.RegisterFilter(HANDLER_ID, isEventHandled)
}

func isEventHandled(event *core.Event, params map[string]string) bool {
	searchText := params["contains"]
	eventBody := event.Data.(github.IssueCommentPayload).Comment.Body
	ruleApplies := strings.Contains(strings.ToUpper(searchText), strings.ToUpper(eventBody))
	return ruleApplies
}
