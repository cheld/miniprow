package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cheld/miniprow/pkg/common/core"
	"github.com/cheld/miniprow/pkg/hook/rules/handlers"
)

const (
	HANDLER_ID    = "http"
	PARAM_URL     = "url"
	PARAM_METHOD  = "method"
	PARAM_BODY    = "body"
	PARAM_HEADERS = "headers"

	VALUE_POST = "post"
)

func init() {
	handlers.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(event *core.Event, params map[string]interface{}) {
	url := params[PARAM_URL].(string)
	event.Log("URL: %v", url)
	method := params[PARAM_METHOD].(string)
	body := params[PARAM_BODY].(string)
	event.Log("BODY: %v", body)
	headers := params[PARAM_HEADERS].(map[string]string)
	event.Log("HEADERS: %v", headers)

	// create http request
	httpClient := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		event.Err("http request failed for action %v, error: %v", HANDLER_ID, err)
		return
	}
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	// execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		event.Err("http request failed for action %v, error: %v", HANDLER_ID, err)
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		event.Err("http request failed for action %v, error: %v", HANDLER_ID, err)
		return
	}
	event.Log("Github response: %v", string(respBody))
}
