package filters

import "github.com/cheld/miniprow/pkg/common/core"

var (
	filter = map[string]Filter{}
)

// TriggerHandler defines the function contract for all triggers.
type Filter func(*core.Event, map[string]string) bool

func RegisterFilter(name string, fn Filter) {
	filter[name] = fn
}

func GetFilter(name string) Filter {
	return filter[name]
}
