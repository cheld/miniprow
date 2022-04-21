package http

import (
	"github.com/cheld/miniprow/pkg/piper/actions"
	"github.com/cheld/miniprow/pkg/piper/config"
)

const (
	actionName = "http"
)

func init() {
	actions.RegisterHandler(actionName, handleAction)
}

func handleAction(rule config.Rule, event config.Event) {

}
