package http

import (
	"github.com/cheld/miniprow/pkg/common/core"
	"github.com/cheld/miniprow/pkg/hook/rules/filters"
)

const (
	HANDLER_ID = "http_request"
)

func init() {
	filters.RegisterFilter(HANDLER_ID, isEventHandled)
}

func isEventHandled(event *core.Event, params map[string]string) bool {

	return true
}
