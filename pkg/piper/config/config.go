package config

import (
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
	Log  []string
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
	Tigger string
	When   map[string]string
}

type Action struct {
	Action string
	With   map[string]interface{}
}

func (config *Configuration) Filter(event Event) []Rule {
	matchingRules := []Rule{}
	for _, rule := range config.Rules {
		if rule.If.Tigger == event.Type {
			resolved, _ := util.ProcessAllTemplates(rule, event)
			matchingRules = append(matchingRules, resolved.(Rule))
		}
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
