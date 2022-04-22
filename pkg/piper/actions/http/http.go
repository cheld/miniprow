package http

import (
	"github.com/cheld/miniprow/pkg/piper/actions"
	"github.com/cheld/miniprow/pkg/piper/config"
)

const (
	HANDLER_ID   = "http"
	PARAM_URL    = "url"
	PARAM_METHOD = "method"
	PARAM_BODY   = "body"

	VALUE_POST = "post"
)

func init() {
	actions.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(params map[string]string, event config.Event) {

}
