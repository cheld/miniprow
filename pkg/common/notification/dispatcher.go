package notification

import (
	"context"
	"fmt"

	"github.com/cheld/miniprow/pkg/common/core"
)

type Dispatcher struct {
	events map[string]Listener
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		events: make(map[string]Listener),
	}
}

func (d *Dispatcher) Register(listener Listener, names ...string) error {
	for _, name := range names {
		d.events[name] = listener
	}
	return nil
}

func (d *Dispatcher) Dispatch(event *core.Event, tentant core.Tenant) error {
	if _, ok := d.events[event.Type]; !ok {
		return fmt.Errorf("no listener registered for even '%s'", event.Type)
	}
	handler := d.events[event.Type]
	handler(event, tentant, context.Background())
	return nil
}
