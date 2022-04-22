package config

import (
	"encoding/json"
	"fmt"

	"github.com/cheld/miniprow/pkg/common/util"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Rules []Rule
}

type Event struct {
	Type string
	Data interface{}
	logs []string
}

func (event *Event) Log(message string, params ...interface{}) {
	line := fmt.Sprintf(message, params...)
	event.logs = append(event.logs, line)
}

func (event *Event) Trail() string {
	s, _ := json.MarshalIndent(event.logs, "", "\t")
	return string(s)
}

type Tenant struct {
	Config  Configuration
	Environ map[string]string
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

func (config *Configuration) Filter(event *Event) []Rule {
	matchingRules := []Rule{}
	for _, rule := range config.Rules {
		event.Log(rule.If.Trigger)
		if rule.If.Trigger == event.Type {
			resolved, _ := util.ProcessAllTemplates(rule, event)
			matchingRules = append(matchingRules, resolved.(Rule))
			event.Log("Matching rule %v", rule.If.Trigger)
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
