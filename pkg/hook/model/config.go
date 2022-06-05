package model

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/common/core"
	"github.com/cheld/miniprow/pkg/common/util"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Rules []Rule
}

type Rule struct {
	If   Trigger
	Then Action
}

type Trigger struct {
	Trigger string
	When    map[string]string
}

type Action struct {
	Action string
	With   map[string]interface{}
}

func (config *Configuration) Filter(event *core.Event) []Rule {
	matchingRules := []Rule{}
	for _, rule := range config.Rules {
		if rule.If.Trigger == event.Type {
			resolved, err := util.ProcessAllTemplates(rule, event)
			if err != nil {
				event.Err("Template for rule %v cannot be processed: %v", rule.If.Trigger, err)
				return matchingRules
			}
			matchingRules = append(matchingRules, resolved.(Rule))
			event.Log("Rule configuration %v found", rule.If.Trigger)
		}
	}
	if len(matchingRules) == 0 {
		event.Log("No matching rule found")
	}
	return matchingRules
}

func Load(cfg *[]byte) (Configuration, error) {
	var yamlConfig Configuration
	err := yaml.Unmarshal(*cfg, &yamlConfig)
	if err != nil {
		return Configuration{}, fmt.Errorf("Error parsing YAML file: %s", err)
	}
	return yamlConfig, nil
}
