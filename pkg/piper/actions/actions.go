package actions

import "github.com/cheld/miniprow/pkg/piper/config"

var (
	handlers = map[string]ActionHandler{}
)

// TriggerHandler defines the function contract for all triggers.
type ActionHandler func(map[string]string, config.Event)

func RegisterHandler(name string, fn ActionHandler) {
	handlers[name] = fn
}

func GetHandler(name string) ActionHandler {
	return handlers[name]
}

func Handle(triggeredRules []config.Rule, event config.Event, tenant config.Tenant) {
	rules := tenant.Config.Filter(event.Type)
	for _, rule := range rules {
		handler := handlers[rule.Then.Action]
		handler(rule.Then.With, event)
	}
}
