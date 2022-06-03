package github

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/common/util"
	"github.com/cheld/miniprow/pkg/hook/model"
	filter "github.com/cheld/miniprow/pkg/hook/rules/filters/github"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers/http"
	"github.com/go-playground/webhooks/v6/github"
)

const (
	HANDLER_ID    = "github_comment"
	PARAM_COMMENT = "comment"
)

func init() {
	handlers.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(event *model.Event, params map[string]interface{}) {
	if _, ok := event.Data.(github.IssueCommentPayload); !ok {
		event.Err("Action %v can only be combined with Trigger %v", HANDLER_ID, filter.HANDLER_ID)
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
	httptHandler := handlers.GetHandler(http.HANDLER_ID)
	httptHandler(event, params)
}
