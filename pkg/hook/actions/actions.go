package actions

import "github.com/cheld/miniprow/pkg/hook/config"

var (
	handlers = map[string]ActionHandler{}
)

// TriggerHandler defines the function contract for all triggers.
type ActionHandler func(map[string]interface{}, config.Event)

func RegisterHandler(name string, fn ActionHandler) {
	handlers[name] = fn
}

func GetHandler(name string) ActionHandler {
	return handlers[name]
}

func Handle(triggeredRules []config.Rule, event *config.Event, tenant config.Tenant) {
	for _, rule := range triggeredRules {
		handler := handlers[rule.Then.Action]
		handler(rule.Then.With, *event)
	}
}
