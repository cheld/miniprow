package triggers

import "github.com/cheld/miniprow/pkg/piper/config"

var (
	handlers = map[string]TriggerHandler{}
)

// TriggerHandler defines the function contract for all triggers.
type TriggerHandler func(config.Event, config.Rule) bool

func RegisterHandler(name string, fn TriggerHandler) {
	handlers[name] = fn
}

func Handle(event config.Event, tenant config.Tenant) []config.Rule {
	handler := handlers[event.Type]
	rules := tenant.Config.Filter(event)
	triggeredRules := []config.Rule{}
	for _, rule := range rules {
		if handler(event, rule) {
			triggeredRules = append(triggeredRules, rule)
		}
	}
	return triggeredRules
}
