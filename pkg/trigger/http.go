package trigger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cheld/cicd-bot/pkg/config"
)

func ExecuteHttp(trigger *config.Trigger, task config.Task) {
	fmt.Println("Executing http trigger.")
	s, _ := config.ProcessAllTemplates(trigger.Spec, task)
	spec := s.(map[string]interface{})
	url := spec["url"].(string)
	method := spec["method"].(string)
	data := spec["data"].(string)

	httpClient := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
	}
	token := task.Environ["SECRET_GITHUB"]
	fmt.Println(task.Environ)
	fmt.Println(token)
	req.Header.Add("Authorization", "token "+token)
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
