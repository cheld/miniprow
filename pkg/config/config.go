package config

import (
	"fmt"
	"io/ioutil"
	"log"
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
	Events   []Event
	Triggers []Trigger
}

func (config *Configuration) FindMatchingEvent(source, eventType string, eventInput EventInput) *Event {
	for _, event := range config.Events {
		if strings.EqualFold(event.Source, source) &&
			strings.EqualFold(event.Type, eventType) &&
			event.IsMatching(eventInput) {
			return &event
		}
	}
	return nil
}

//func (config *Configuration) getTrigger(name string) Trigger {
//	for _, trigger := range config.Triggers {
//		if trigger.Name == name {
//			return trigger
//		}
//	}
//	return Trigger{}
//}

type EventInput struct {
	Objectiv string
	Input    map[string]interface{}
}

type Event struct {
	Source      string
	Type        string
	If_contains string
	If_equals   string
	If_true     string
	Trigger     string
	Values      map[string]interface{}
}

func (event *Event) IsMatching(eventInput EventInput) bool {
	contains := true
	if event.If_contains != "" {
		contains = strings.Contains(strings.ToUpper(eventInput.Objectiv), strings.ToUpper(event.If_contains))
	}
	equals := true
	if event.If_equals != "" {
		equals = strings.EqualFold(eventInput.Objectiv, event.If_equals)
	}
	condition := true
	if event.If_true != "" {
		result, err := ProcessTemplate(event.If_true, eventInput)
		if err != nil {
			result = "false"
			log.Println(err)
		}
		condition, _ = strconv.ParseBool(result)
	}
	return contains && equals && condition
}

func (event *Event) Handle(eventInput EventInput) TriggerInput {
	triggerInput := TriggerInput{}
	triggerInput.Name = event.Trigger
	result, err := ProcessAllTemplates(event.Values, eventInput)
	if err != nil {
		panic(err)
	}
	triggerInput.Values = result.(map[string]interface{})
	return triggerInput
}

type TriggerInput struct {
	Name   string
	Values map[string]interface{}
}

type Trigger struct {
	Name      string
	Type      string
	Arguments map[string]interface{}
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
