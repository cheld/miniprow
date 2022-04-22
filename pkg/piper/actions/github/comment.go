package github

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/piper/actions"
	"github.com/cheld/miniprow/pkg/piper/actions/http"
	"github.com/cheld/miniprow/pkg/piper/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	HANDLER_ID    = "github_comment"
	PARAM_COMMENT = "comment"
)

func init() {
	actions.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(params map[string]string, event config.Event) {
	issueNumber := event.Data.(github.IssueCommentPayload).Issue.Number
	comment := params[PARAM_COMMENT]

	//headers:
	//Authorization: 'token {{.Environ.SECRET_GITHUB }}'

	params[http.PARAM_URL] = fmt.Sprintf("https://api.github.com/repos/cheld/code-snippets/issues/%v/comments", issueNumber)
	params[http.PARAM_METHOD] = http.VALUE_POST
	params[http.PARAM_BODY] = fmt.Sprintf("{ \"body\": \"%v\" }", comment)
	httptHandler := actions.GetHandler(http.HANDLER_ID)
	httptHandler(params, event)

}
