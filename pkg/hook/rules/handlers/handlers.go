package handlers

import config "github.com/cheld/miniprow/pkg/hook/model"

var (
	handlers = map[string]Handler{}
)

// TriggerHandler defines the function contract for all triggers.
type Handler func(*config.Event, map[string]interface{})

func RegisterHandler(name string, fn Handler) {
	handlers[name] = fn
}

func GetHandler(name string) Handler {
	return handlers[name]
}
