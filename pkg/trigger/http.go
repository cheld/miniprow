package trigger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cheld/cicd-bot/pkg/config"
)

func ExecuteHttp(trigger config.Trigger, task config.Task) {
	s, _ := config.ProcessAllTemplates(trigger.Spec, task)
	spec := s.(map[string]string)
	url := spec["url"]
	method := spec["method"]
	data := spec["data"]

	httpClient := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "token 4accfefc3f1acb3896748654c0148ab352826cf3")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}
