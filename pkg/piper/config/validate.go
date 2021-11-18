package config

import (
	"errors"
	"fmt"
	"strings"
)

func Validate(cfg Configuration) error {
	if len(cfg.Rules) == 0 {
		return errors.New("No rules defined in configuration")
	}
	for i, rule := range cfg.Rules {
		if rule.Event == "" {
			return fmt.Errorf("Rule[%d] has no event defined", i)
		}
		if rule.If_contains == "" && rule.If_equals == "" && rule.If_true == "" {
			return fmt.Errorf("Event (%s) must define a matching rule", rule.Event)
		}
		if len(rule.Then) == 0 {
			return fmt.Errorf("Event (%s) has no trigger definition", rule.Event)
		}
		if len(rule.Then) > 0 && cfg.GetTrigger(rule.Then[0].Apply) == nil {
			return fmt.Errorf("Event (%s) references a trigger that does not exist", rule.Event)
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
