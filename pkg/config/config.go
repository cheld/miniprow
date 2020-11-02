package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	//"gopkg.in/yaml.v2"
	//sigs.k8s.io/yaml"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Serve struct {
		Secret string
		Port   int
	}
	Events  Event
	Trigger Trigger
}

//func (config *Configuration) getTrigger(name string) Trigger {
//	for _, trigger := range config.Triggers {
//		if trigger.Name == name {
//			return trigger
//		}
//	}
//	return Trigger{}
//}

type Event struct {
	Github []EventSpec
	Cli    []EventSpec
}

type InputData struct {
	Objectiv string
	Input    map[string]interface{}
}

type EventData struct {
	Name   string
	Values map[string]interface{}
}

type EventSpec struct {
	Type     string
	Contains string
	Equals   string
	Rule     string
	Trigger  string
	Values   map[string]interface{}
}

func (rule *EventSpec) IsMatching(inputData InputData) bool {
	contains := true
	if rule.Contains != "" {
		contains = strings.Contains(inputData.Objectiv, rule.Contains)
	}
	equals := true
	if rule.Equals != "" {
		equals = inputData.Objectiv == rule.Equals
	}
	condition := true
	if rule.Rule != "" {
		result := ProcessTemplate(rule.Rule, inputData)
		condition, _ = strconv.ParseBool(result)
	}
	return contains && equals && condition
}

func (spec *EventSpec) Process(inputData InputData) EventData {
	eventData := EventData{}
	eventData.Name = spec.Trigger
	eventData.Values = ProcessAllTemplates(spec.Values, inputData).(map[string]interface{})
	return eventData
}

type Trigger struct {
	Http  []HttpTrigger
	Debug []DebugTrigger
}

type HttpTrigger struct {
	Name string
	Url  string
}

type DebugTrigger struct {
	Name string
	Text string
}

func Load(filename string) Configuration {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
	}

	var yamlConfig Configuration
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	return yamlConfig
}
