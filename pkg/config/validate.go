package config

import (
	"errors"
	"fmt"
	"strings"
)

func Validate(cfg Configuration) error {
	if len(cfg.Events) == 0 {
		return errors.New("No events defined in configuration")
	}
	for i, event := range cfg.Events {
		if event.Source == "" {
			return fmt.Errorf("Event[%d] has no source defined", i)
		}
		if event.Type == "" {
			return fmt.Errorf("Event[%d] has no type defined", i)
		}
		if event.If_contains == "" && event.If_equals == "" && event.If_true == "" {
			return fmt.Errorf("Event (%s,%s) must define a matching rule", event.Source, event.Type)
		}
		if event.Trigger == "" {
			return fmt.Errorf("Event (%s,%s) has no trigger definition", event.Source, event.Type)
		}
		if event.Trigger != "" && cfg.Trigger(event.Trigger) == nil {
			return fmt.Errorf("Event (%s,%s) references a trigger that does not exist", event.Source, event.Type)
		}

	}
	if len(cfg.Triggers) == 0 {
		return errors.New("No triggers defined in configuration")
	}
	for i, trigger := range cfg.Triggers {
		if trigger.Name == "" {
			return fmt.Errorf("Trigger[%d] has no name defined", i)
		}
		if trigger.Type == "" {
			return fmt.Errorf("Trigger '%s' has no type defined", trigger.Name)
		}
		if trigger.Type != "" && strings.ToLower(trigger.Type) != "http" && strings.ToLower(trigger.Type) != "debug" {
			return fmt.Errorf("Type of trigger '%s' must be either http or debug", trigger.Name)
		}
	}
	return nil
}
