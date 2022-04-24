package github

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/common/util"
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/plugins/actions"
	"github.com/cheld/miniprow/pkg/hook/plugins/actions/http"
	trigger "github.com/cheld/miniprow/pkg/hook/plugins/triggers/github"
	"github.com/go-playground/webhooks/v6/github"
)

const (
	HANDLER_ID    = "github_comment"
	PARAM_COMMENT = "comment"
)

func init() {
	actions.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(params map[string]interface{}, event *model.Event) {
	if _, ok := event.Data.(github.IssueCommentPayload); !ok {
		event.Err("Action %v can only be combined with Trigger %v", HANDLER_ID, trigger.HANDLER_ID)
		return
	}
	issueNumber := event.Data.(github.IssueCommentPayload).Issue.Number
	event.Log("Commenting on Github issue %d", issueNumber)
	comment := params[PARAM_COMMENT]
	event.Log("Comment: %v", comment)

	headers := map[string]string{}
	headers["Authorization"] = "token " + util.Environment.Value("GITHUB_TOKEN").String()
	headers["Content-Type"] = "application/json"

	params[http.PARAM_URL] = fmt.Sprintf("https://api.github.com/repos/cheld/testpython/issues/%v/comments", issueNumber)
	params[http.PARAM_METHOD] = "POST" //http.VALUE_POST
	params[http.PARAM_BODY] = fmt.Sprintf("{\"body\": \"![cat](https://static.elle.de/1200x630/smart/images/2016-03/15325082_523f72f18b.jpg)\"}")
	params[http.PARAM_HEADERS] = headers
	httptHandler := actions.GetHandler(http.HANDLER_ID)
	httptHandler(params, event)
}
