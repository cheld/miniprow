package rules

import (
	"github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/rules/filters"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers"
)

type listener struct {
	eventFilter  filters.Filter
	filterParam  map[string]string
	eventHandler handlers.Handler
	handlerParam map[string]interface{}
}

func NewRuleBasedListeners(rules []model.Rule) []listener {
	listeners := []listener{}
	for _, rule := range rules {
		listeners = append(listeners, NewListener(rule))
	}
	return listeners
}

func NewListener(rule model.Rule) listener {
	l := listener{}
	l.eventFilter = filters.GetFilter(rule.If.Trigger)
	l.eventHandler = handlers.GetHandler(rule.Then.Action)
	l.filterParam = rule.If.When
	l.handlerParam = rule.Then.With
	return l
}

func (l *listener) Handle(event model.Event) {
	if l.eventFilter(event, l.filterParam) {
		return
	}
	l.eventHandler(&event, l.handlerParam)

}
