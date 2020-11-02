package trigger

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

func Stdout(trigger config.DebugTrigger, eventData config.EventData) {
	t := config.ProcessTemplate(trigger.Text, eventData)
	fmt.Println(t)
}
