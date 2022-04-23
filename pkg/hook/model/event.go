package model

import (
	"encoding/json"
	"fmt"
)

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
