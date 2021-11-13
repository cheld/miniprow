package event

import (
	"encoding/json"

	"github.com/cheld/miniprow/pkg/common/util"
	"github.com/cheld/miniprow/pkg/piper/config"
	"github.com/golang/glog"
)

func (handler *Handler) HandleHttp(body []byte, path string) []config.Task {

	// parse body
	var payload interface{}
	if len(body) > 0 {
		err := json.Unmarshal(body, &payload)
		if err != nil {
			glog.Errorf("Not possible to parse request body %s, %v", string(body), err)
			return []config.Task{}
		}
	} else {
		payload = ""
	}
	request := config.Request{
		Value:   string(body),
		Payload: payload,
	}
	ctx := config.Ctx{
		Request: request,
		Environ: *util.Environment.Map(),
	}

	// handle event
	event := handler.config.GetMatchingRule("http", ctx)
	if event == nil {
		glog.V(5).Infof("No event found for value %s\n", ctx.Request.Value)
		return []config.Task{}
	}

	// build execution task
	task, err := event.BuildTask(ctx)
	if err != nil {
		glog.Errorf("Cannot handle event: %v. Error: %v", event.Trigger, err)
		return []config.Task{}
	}
	return []config.Task{task}
}
