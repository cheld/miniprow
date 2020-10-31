package config

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	//"gopkg.in/yaml.v2"
	//sigs.k8s.io/yaml"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Serve struct {
		Secret string
		Port   int
	}
	Events       Event
	Destinations Destination
	Cron         Cron
}

//func (config *Configuration) getTrigger(name string) Trigger {
//	for _, trigger := range config.Triggers {
//		if trigger.Name == name {
//			return trigger
//		}
//	}
//	return Trigger{}
//}

type Cron struct {
}

type Event struct {
	Github GithubEvent
	Cli    CliEvent
}

type GithubEvent struct {
	Comment []Rule
}

type CliEvent struct {
	Stdin []Rule
}

type EventCtx struct {
	Objectiv string
	Input    map[string]interface{}
}

type DestinationCtx struct {
	Name   string
	Values map[string]interface{}
}

type Rule struct {
	Contains    string
	Equals      string
	Condition   string
	Destination string
	Values      map[string]interface{}
}

func (rule *Rule) IsMatching(ctx EventCtx) bool {
	contains := true
	if rule.Contains != "" {
		contains = strings.Contains(ctx.Objectiv, rule.Contains)
	}
	equals := true
	if rule.Equals != "" {
		equals = ctx.Objectiv == rule.Contains
	}
	condition := true
	if rule.Condition != "" {
		var tpl bytes.Buffer
		t, _ := template.New("Condition").Parse(rule.Condition)
		_ = t.Execute(&tpl, ctx)
		result := tpl.String()
		condition, _ = strconv.ParseBool(result)
	}
	return contains && equals && condition
}

func (rule *Rule) DestinationCtx(input EventCtx) DestinationCtx {
	destination := DestinationCtx{}
	destination.Name = rule.Destination
	destination.Values = input.Copy(rule.Values).(map[string]interface{})
	return destination
}

type Destination struct {
	Http  HttpDestination
	Debug DebugDestination
}

type HttpDestination struct {
	Get []HttpGetDestination
}

type HttpGetDestination struct {
	Name string
	Url  string
}

type DebugDestination struct {
	Info []DebugStdoutDestination
}

type DebugStdoutDestination struct {
	Name string
	Text string
}

//func (trigger *Trigger) httpParamer(target string) HttpParameter {
//	parameter := HttpParameter{}
//	input := struct {
//		Target string
//	}{
//		target,
//	}
//
//	var tpl bytes.Buffer
//	t, _ := template.New("todos").Parse(trigger.Url)
//	_ = t.Execute(&tpl, input)
//
//	parameter.Url = tpl.String()
//	return parameter
//}

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

	//Sfmt.Printf("Result: %v\n", yamlConfig)
	return yamlConfig
}

// Copy creates a deep copy of whatever is passed to it and returns the copy
// in an interface{}.  The returned value will need to be asserted to the
// correct type.
func (ctx *EventCtx) Copy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	ctx.copyRecursive(original, cpy)

	// Return the copy as an interface.
	return cpy.Interface()
}

// copyRecursive does the actual copying of the interface. It currently has
// limited support for what it can handle. Add as needed.
func (ctx *EventCtx) copyRecursive(original, cpy reflect.Value) {

	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		ctx.copyRecursive(originalValue, cpy.Elem())

	case reflect.Interface:
		// If this is a nil, don't do anything
		if original.IsNil() {
			return
		}
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		ctx.copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)

	case reflect.Struct:
		t, ok := original.Interface().(time.Time)
		if ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}
		// Go through each field of the struct and copy it.
		for i := 0; i < original.NumField(); i++ {
			// The Type's StructField for a given field is checked to see if StructField.PkgPath
			// is set to determine if the field is exported or not because CanSet() returns false
			// for settable fields.  I'm not sure why.  -mohae
			if original.Type().Field(i).PkgPath != "" {
				continue
			}
			ctx.copyRecursive(original.Field(i), cpy.Field(i))
		}

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			ctx.copyRecursive(original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			ctx.copyRecursive(originalValue, copyValue)
			copyKey := ctx.Copy(key.Interface())
			cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		if original.Kind() == reflect.String {
			fmt.Println("test2------------")
			fmt.Println(original.String())
			n := ExecuteTemplate(original.String(), ctx)
			newValue := reflect.ValueOf(n)
			cpy.Set(newValue)
		} else {
			cpy.Set(original)
		}
	}
}

func ExecuteTemplate(tpl string, data interface{}) string {
	var result bytes.Buffer
	t, _ := template.New("tmp").Parse(tpl)
	_ = t.Execute(&result, data)
	return result.String()
}
