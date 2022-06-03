package misc

import (
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers/github"
)

const (
	actionName = "cat"
)

func init() {
	handlers.RegisterHandler(actionName, handleAction)
}

func handleAction(event *model.Event, params map[string]interface{}) {
	params[github.PARAM_COMMENT] = "here is the cat"
	commentHandler := handlers.GetHandler(github.HANDLER_ID)
	commentHandler(event, params)
}
