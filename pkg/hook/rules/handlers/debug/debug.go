package debug

import (
	"github.com/cheld/miniprow/pkg/common/core"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers"
)

const (
	HANDLER_ID = "debug"
)

func init() {
	handlers.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(event *core.Event, params map[string]interface{}) {
	event.Log("Debug action executed")
}
