package filters

import (
	"github.com/cheld/miniprow/pkg/hook/model"
)

var (
	filter = map[string]Filter{}
)

// TriggerHandler defines the function contract for all triggers.
type Filter func(model.Event, map[string]string) bool

func RegisterFilter(name string, fn Filter) {
	filter[name] = fn
}

func GetFilter(name string) Filter {
	return filter[name]
}
