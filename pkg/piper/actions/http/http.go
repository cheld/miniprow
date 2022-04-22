package http

import (
	"bytes"
	"net/http"

	"github.com/cheld/miniprow/pkg/piper/actions"
	"github.com/cheld/miniprow/pkg/piper/config"
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
	actions.RegisterHandler(HANDLER_ID, handleAction)
}

func handleAction(params map[string]interface{}, event config.Event) {

	url := params[PARAM_URL].(string)
	method := params[PARAM_METHOD].(string)
	body := params[PARAM_BODY].(string)
	headers := params[PARAM_HEADERS].(map[string]string)

	// create http request
	httpClient := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBufferString(body))
	//if err != nil {
	//	return fmt.Errorf("Error occured when creating http trigger '%s'. %s", trigger.Name, err)
	//}
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	// execute request
	resp, _ := httpClient.Do(req)
	//if err != nil {
	//	return fmt.Errorf("Error occured when executing trigger %s, with error %s", trigger.Name, err)
	//}
	defer resp.Body.Close()
	//respBody, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return fmt.Errorf("Trigger %s failed, error %s", trigger.Name, err)
	//}
}
