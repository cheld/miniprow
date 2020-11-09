package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

type Source struct {
	Value   string
	Payload interface{}
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

func (config *Configuration) Event(source, eventType string, sourceData Source) *Event {
	for _, event := range config.Events {
		if strings.EqualFold(event.Source, source) &&
			strings.EqualFold(event.Type, eventType) &&
			event.IsMatching(sourceData) {
			return &event
		}
	}
	return nil
}

func (event *Event) IsMatching(source Source) bool {
	contains := true
	if event.If_contains != "" {
		contains = strings.Contains(strings.ToUpper(source.Value), strings.ToUpper(event.If_contains))
	}
	equals := true
	if event.If_equals != "" {
		equals = strings.EqualFold(source.Value, event.If_equals)
	}
	condition := true
	if event.If_true != "" {
		result, err := ProcessTemplate(event.If_true, source)
		if err != nil {
			result = "false"
			log.Println(err)
		}
		condition, _ = strconv.ParseBool(result)
	}
	return contains && equals && condition
}

func (event *Event) NewTask(source Source) (Task, error) {
	task := Task{}
	task.Trigger = event.Trigger
	result, err := ProcessAllTemplates(event.Values, source)
	if err != nil {
		return task, fmt.Errorf("Cannot process: %v. Error: %v", task.Trigger, err)
	}
	task.Values = result.(map[string]interface{})
	task.Env = Env()
	return task, nil
}

type Task struct {
	Trigger string
	Values  map[string]interface{}
	Env     map[string]string
}

type Trigger struct {
	Name string
	Type string
	Spec map[string]interface{}
}

func (config *Configuration) Trigger(name string) *Trigger {
	for _, trigger := range config.Triggers {
		if strings.EqualFold(trigger.Name, name) {
			return &trigger
		}
	}
	return nil
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

func Env() map[string]string {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		keyValue := strings.Split(entry, "=")
		env[keyValue[0]] = keyValue[1]
	}
	return env
}
