package trigger

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

func ExecuteDebug(trigger config.Trigger, task config.Task) {
	t, _ := config.ProcessAllTemplates(trigger.Spec, task)
	values := t.(map[string]interface{})
	message := values["stdout"]
	fmt.Println(message)
}
