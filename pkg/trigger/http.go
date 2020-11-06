package trigger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cheld/cicd-bot/pkg/config"
)

func ExecuteHttp(trigger config.Trigger, triggerInput config.TriggerInput) {
	url := trigger.Arguments["url"].(string)
	method := trigger.Arguments["method"].(string)
	data := trigger.Arguments["data"].(string)

	u, _ := config.ProcessTemplate(url, triggerInput)
	httpClient := &http.Client{}

	req, err := http.NewRequest(method, u, bytes.NewBufferString(data))
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
