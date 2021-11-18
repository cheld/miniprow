package trigger

import (
	"fmt"

	"github.com/cheld/miniprow/pkg/piper/config"
)

func ExecuteDebug(trigger *config.Trigger, ctx config.Ctx) error {
	t, _ := config.ProcessAllTemplates(trigger.Spec, ctx)
	values := t.(map[string]interface{})
	message := values["stdout"]
	fmt.Println(message)
	return nil
}
