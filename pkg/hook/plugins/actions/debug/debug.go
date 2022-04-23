package debug

import (
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/plugins/actions"
)

const (
	HANDLER_ID = "debug"
)

func init() {
	actions.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(params map[string]interface{}, event *model.Event) {
	event.Log("Debug action executed")
}
