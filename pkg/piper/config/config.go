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
	Request Request
	Environ map[string]string
	Trigger map[string]string
}

type Request struct {
	Event   string
	Value   string
	Payload interface{}
}

type Rule struct {
	Event       string
	If_contains string
	If_equals   string
	If_true     string
	Trigger     map[string]interface{}
}

func (config *Configuration) GetMatchingRule(event string, sourceData Ctx) *Rule {
	for _, rule := range config.Rules {
		if strings.EqualFold(rule.Event, event) &&
			rule.IsMatching(sourceData) {
			return &rule
		}
	}
	return nil
}

func (rule *Rule) IsMatching(ctx Ctx) bool {
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

func (event *Rule) BuildTask(ctx Ctx) (Task, error) {
	task := Task{}
	task.Trigger = event.Trigger["action"].(string)
	result, err := ProcessAllTemplates(event, ctx)
	if err != nil {
		return task, fmt.Errorf("Cannot process: %v. Error: %v", task.Trigger, err)
	}
	task.Values = result.(map[string]interface{})
	task.Environ = ctx.Environ
	return task, nil
}

type Task struct {
	Trigger string
	Values  map[string]interface{}
	Environ map[string]string
}

type Trigger struct {
	Name string
	Type string
	Spec map[string]interface{}
}

func (config *Configuration) FindTrigger(name string) *Trigger {
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
