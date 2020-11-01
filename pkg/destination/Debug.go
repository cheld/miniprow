package destination

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cheld/cicd-bot/pkg/config"
)

func Stdout(destination config.DebugStdout, eventData config.EventData) {
	t := ExecuteTemplate(destination.Text, eventData)
	fmt.Println(t)
}

func ExecuteTemplate(tpl string, data interface{}) string {
	var result bytes.Buffer
	t, _ := template.New("tmp").Parse(tpl)
	_ = t.Execute(&result, data)
	return result.String()
}
