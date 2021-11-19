package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Serve struct {
		Secret string
		Port   int
	}
	Rules    []Rule
	Triggers []Trigger
}

type Ctx struct {
	Request RequestCtx
	Environ map[string]string
	Trigger map[string]TriggerCtx
}

type RequestCtx struct {
	Event   string
	Value   string
	Payload interface{}
}

type TriggerCtx struct {
	Input map[string]string
}

type Rule struct {
	Event       string
	If_contains string
	If_equals   string
	If_true     string
	Then        []Then
}

type Then struct {
	Apply string
	With  map[string]string
}

func (config *Configuration) GetFirstMatchingRule(ctx *Ctx) *Rule {
	for _, rule := range config.Rules {
		if rule.IsMatching(ctx) {
			return &rule
		}
	}
	return nil
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

type Trigger struct {
	Name string
	Type string
	Spec map[string]interface{}
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
