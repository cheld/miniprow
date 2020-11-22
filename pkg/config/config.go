package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

const (
	cfgDefaultName = "cicd-bot.yaml"
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
	Environ map[string]string
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

func (config *Configuration) FindEvent(source, eventType string, sourceData Source) *Event {
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

func (event *Event) BuildTask(source Source) (Task, error) {
	task := Task{}
	task.Trigger = event.Trigger
	result, err := ProcessAllTemplates(event.Values, source)
	if err != nil {
		return task, fmt.Errorf("Cannot process: %v. Error: %v", task.Trigger, err)
	}
	task.Values = result.(map[string]interface{})
	task.Environ = source.Environ
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

func Load(filename string) (Configuration, error) {
	yamlFile, err := readFile(filename)
	if err != nil {
		return Configuration{}, fmt.Errorf("Error reading YAML file: %s", err)
	}
	var yamlConfig Configuration
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return Configuration{}, fmt.Errorf("Error parsing YAML file: %s", err)
	}
	err = Validate(yamlConfig)
	if err != nil {
		return Configuration{}, fmt.Errorf("Error validating YAML file: %s", err)
	}
	return yamlConfig, nil
}

func readFile(filename string) ([]byte, error) {
	var yamlFile []byte
	if filename != "" {
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			return yamlFile, fmt.Errorf("Error reading YAML file: %s", err)
		}
		return yamlFile, nil
	}
	etcPath := fmt.Sprintf("/etc/.%s", cfgDefaultName)
	yamlFile, err := ioutil.ReadFile(etcPath)
	if err == nil {
		return yamlFile, nil
	}
	home, _ := homedir.Dir()
	homepath := fmt.Sprintf("%s/.%s", home, cfgDefaultName)
	yamlFile, err = ioutil.ReadFile(homepath)
	if err == nil {
		return yamlFile, nil
	}
	yamlFile, err = ioutil.ReadFile(cfgDefaultName)
	if err == nil {
		return yamlFile, nil
	}
	return yamlFile, fmt.Errorf("Cound not reading config file from %s, %s, %s", etcPath, homepath, cfgDefaultName)

}
