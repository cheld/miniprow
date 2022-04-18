package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Rules []Rule
}

func (config *Configuration) ProcessAllTemplates(ctx *Ctx) (Configuration, error) {
	result, err := ProcessAllTemplates(config, ctx)
	return result.(Configuration), err
}

func (config *Configuration) GetFirstMatchingRule(ctx *Ctx) *Rule {
	for _, rule := range config.Rules {
		if rule.IsMatching(ctx) {
			return &rule
		}
	}
	return nil
}

type Event struct {
	Type    string
	Data    interface{}
	Rule    Rule
	Environ map[string]string
	Logs    []string
}

type Rule struct {
	If   Condition
	Then Task
}

type Condition struct {
	Tigger string
	When   map[string]string
}

type Task struct {
	Action string
	With   map[string]string
}

func (rule *Rule) IsMatching(ctx *Ctx) bool {
	if !strings.EqualFold(rule.Event, ctx.Request.Event) {
		return false
	}
	contains := true
	if rule.If_contains != "" {
		contains = strings.Contains(strings.ToUpper(ctx.Request.Value), strings.ToUpper(rule.If_contains))
	}
	equals := true
	if rule.If_equals != "" {
		equals = strings.EqualFold(ctx.Request.Value, rule.If_equals)
	}
	condition := true
	if rule.If_true != "" {
		result, err := ProcessTemplate(rule.If_true, ctx)
		if err != nil {
			result = "false"
			log.Println(err)
		}
		condition, _ = strconv.ParseBool(result)
	}
	return contains && equals && condition
}

//ProcessRuleTemplates(ctx)

func (config *Configuration) ApplyRule(ctx *Ctx) error {
	rule := config.GetFirstMatchingRule(ctx)
	if rule == nil {
		return fmt.Errorf("No rule defined for %v", rule.Event)
	}
	return rule.Apply(ctx)
}

func (rule *Rule) Apply(ctx *Ctx) error {
	for _, then := range rule.Then {
		result, err := ProcessAllTemplates(then.With, ctx)
		if err != nil {
			return fmt.Errorf("Cannot process: %v. Error: %v", rule.Event, err)
		}
		if ctx.Trigger == nil {
			ctx.Trigger = make(map[string]TriggerCtx)
		}
		triggerCtx := TriggerCtx{
			Input: result.(map[string]string),
		}
		ctx.Trigger[then.Apply] = triggerCtx
	}
	return nil
}

//ProcessTriggerTemplates(ctx)

func (config *Configuration) GetTrigger(name string) *Trigger {
	for _, trigger := range config.Triggers {
		if strings.EqualFold(trigger.Name, name) {
			return &trigger
		}
	}
	return nil
}

func Load(cfg *[]byte) (Configuration, error) {
	var yamlConfig Configuration
	err := yaml.Unmarshal(*cfg, &yamlConfig)
	if err != nil {
		return Configuration{}, fmt.Errorf("Error parsing YAML file: %s", err)
	}
	err = Validate(yamlConfig)
	if err != nil {
		return Configuration{}, fmt.Errorf("Error validating YAML file: %s", err)
	}
	return yamlConfig, nil
}
