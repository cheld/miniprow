package debug

import (
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers"
)

const (
	HANDLER_ID = "debug"
)

func init() {
	handlers.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(event *model.Event, params map[string]interface{}) {
	event.Log("Debug action executed")
}
