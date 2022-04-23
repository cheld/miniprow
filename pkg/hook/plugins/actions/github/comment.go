package github

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/hook/actions"
	"github.com/cheld/miniprow/pkg/hook/actions/http"
	"github.com/cheld/miniprow/pkg/hook/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	HANDLER_ID    = "github_comment"
	PARAM_COMMENT = "comment"
)

func init() {
	actions.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(params map[string]interface{}, event config.Event) {
	issueNumber := event.Data.(github.IssueCommentPayload).Issue.Number
	comment := params[PARAM_COMMENT]

	headers := map[string]string{}
	headers["Authorization"] = "token xxxx"

	params[http.PARAM_URL] = fmt.Sprintf("https://api.github.com/repos/cheld/code-snippets/issues/%v/comments", issueNumber)
	params[http.PARAM_METHOD] = http.VALUE_POST
	params[http.PARAM_BODY] = fmt.Sprintf("{ \"body\": \"%v\" }", comment)
	params[http.PARAM_HEADERS] = headers
	httptHandler := actions.GetHandler(http.HANDLER_ID)
	httptHandler(params, event)
}
