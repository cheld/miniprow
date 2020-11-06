package trigger

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

func ExecuteDebug(trigger config.Trigger, triggerInput config.TriggerInput) {
	messageTemplate := trigger.Arguments["stdout"].(string)
	fmt.Println("--")
	fmt.Println(triggerInput)
	t, _ := config.ProcessTemplate(messageTemplate, triggerInput)
	fmt.Println(t)
}
