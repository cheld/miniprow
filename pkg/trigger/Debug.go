package trigger

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

func Stdout(trigger config.Trigger, triggerInput config.TriggerInput) {
	messageTemplate := trigger.Arguments["text"].(string)
	t, _ := config.ProcessTemplate(messageTemplate, triggerInput)
	fmt.Println(t)
}
