package trigger

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/piper/config"
)

func ExecuteDebug(trigger *config.Trigger, task config.Task) error {
	t, _ := config.ProcessAllTemplates(trigger.Spec, task)
	values := t.(map[string]interface{})
	message := values["stdout"]
	fmt.Println(message)
	return nil
}
