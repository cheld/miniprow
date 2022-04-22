package http

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/hook/config"
	"github.com/cheld/miniprow/pkg/hook/triggers"
)

const (
	triggerName = "http_request"
)

func init() {
	fmt.Println("init")
	triggers.RegisterHandler(triggerName, handleEvent)
}

func handleEvent(event config.Event, rule config.Rule) bool {

	return true
}
