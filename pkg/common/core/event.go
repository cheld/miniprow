package core

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type   string
	Data   interface{}
	status status
}

type status struct {
	logs []string
	err  string
}

func (event *Event) Err(message string, params ...interface{}) {
	txt := fmt.Sprintf("ERROR: "+message, params...)
	event.status.err = txt
	event.status.logs = append(event.status.logs, txt)
}

func (event *Event) Log(message string, params ...interface{}) {
	txt := fmt.Sprintf(message, params...)
	event.status.logs = append(event.status.logs, txt)
}

func (event *Event) Trail() string {
	s, _ := json.MarshalIndent(event.status.logs, "", "\t")
	return string(s)
}
