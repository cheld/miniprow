package http

import (
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/rules/filters"
)

const (
	HANDLER_ID = "http_request"
)

func init() {
	filters.RegisterFilter(HANDLER_ID, isEventHandled)
}

func isEventHandled(event model.Event, params map[string]string) bool {

	return true
}
