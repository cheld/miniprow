package http

import (
	config "github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/plugins/triggers"
)

const (
	HANDLER_ID = "http_request"
)

func init() {
	triggers.RegisterHandler(HANDLER_ID, handleEvent)
}

func handleEvent(event config.Event, rule config.Rule) bool {

	return true
}
