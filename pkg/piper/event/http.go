package event

import (
	"encoding/json"

	"github.com/cheld/miniprow/pkg/common/util"
	"github.com/cheld/miniprow/pkg/piper/config"
	"github.com/golang/glog"
)

func (handler *Handler) HandleHttp(body []byte, path string) config.Ctx {

	// parse body
	var payload interface{}
	if len(body) > 0 {
		err := json.Unmarshal(body, &payload)
		if err != nil {
			glog.Errorf("Not possible to parse request body %s, %v", string(body), err)
			return config.Ctx{}
		}
	} else {
		payload = ""
	}
	request := config.RequestCtx{
		Event:   "http",
		Value:   string(body),
		Payload: payload,
	}
	ctx := config.Ctx{
		Request: request,
		Environ: *util.Environment.Map(),
	}

	// handle event
	handler.config.ApplyRule(&ctx)
	return ctx
}
