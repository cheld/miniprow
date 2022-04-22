package misc

import (
	"github.com/cheld/miniprow/pkg/piper/actions"
	"github.com/cheld/miniprow/pkg/piper/actions/github"
	"github.com/cheld/miniprow/pkg/piper/config"
)

const (
	actionName = "cat"
)

func init() {
	actions.RegisterHandler(actionName, handleAction)
}

func handleAction(params map[string]string, event config.Event) {
	params[github.PARAM_COMMENT] = "here is the cat"
	commentHandler := actions.GetHandler(github.HANDLER_ID)
	commentHandler(params, event)
}
