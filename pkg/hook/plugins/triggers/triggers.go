package triggers

import (
	config "github.com/cheld/miniprow/pkg/hook/model"
)

var (
	handlers = map[string]TriggerHandler{}
)

// TriggerHandler defines the function contract for all triggers.
type TriggerHandler func(config.Event, config.Rule) bool

func RegisterHandler(name string, fn TriggerHandler) {
	handlers[name] = fn
}

func Handle(event *config.Event, tenant config.Tenant) []config.Rule {
	triggeredRules := []config.Rule{}
	handler := handlers[event.Type]
	if handler == nil {
		event.Log("No trigger handler implementation for %v", event.Type)
		return triggeredRules
	}
	rules := tenant.Config.Filter(event)
	for _, rule := range rules {
		if handler(*event, rule) {
			event.Log("Trigger %v fired with %v", rule.If.Trigger, rule.If.When)
			triggeredRules = append(triggeredRules, rule)
		}
	}
	return triggeredRules
}
