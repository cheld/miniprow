package trigger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"

	"github.com/cheld/miniprow/pkg/piper/config"
)

func ExecuteHttp(trigger *config.Trigger, task config.Task) error {
	glog.Infof("Executing http trigger '%s'\n", trigger.Name)

	// Get http parameters
	params, err := config.ProcessAllTemplates(trigger.Spec, task)
	if err != nil {
		return fmt.Errorf("Error occured when processing trigger '%s'. %s", trigger.Name, err)
	}
	url := stringValue(params, "url")
	method := stringValue(params, "method")
	body := stringValue(params, "body")
	headers := params.(map[string]interface{})["headers"].(map[string]interface{})
	//glog.Infof("Http parameters:\nurl=%s\nmethod=%s\nbody=%s\nheaders=%v\n\n", url, method, body, headers) //might contain secrets

	// create http request
	httpClient := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("Error occured when creating http trigger '%s'. %s", trigger.Name, err)
	}
	for key, val := range headers {
		req.Header.Add(key, val.(string))
	}

	// execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error occured when executing trigger %s, with error %s", trigger.Name, err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Trigger %s failed, error %s", trigger.Name, err)
	}
	glog.V(6).Infof("Http response from trigger %s\n%s", trigger.Name, string(respBody))
	return nil
}

func stringValue(parameters interface{}, key string) string {
	parametersMap := parameters.(map[string]interface{})
	return parametersMap[key].(string)
}
