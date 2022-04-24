package misc

import (
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/plugins/actions"
	"github.com/cheld/miniprow/pkg/hook/plugins/actions/github"
)

const (
	actionName = "cat"
)

func init() {
	actions.RegisterHandler(actionName, handleAction)
}

func handleAction(params map[string]interface{}, event *model.Event) {
	params[github.PARAM_COMMENT] = "here is the cat"
	commentHandler := actions.GetHandler(github.HANDLER_ID)
	commentHandler(params, event)
}
