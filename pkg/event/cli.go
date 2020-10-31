package event

import (
	"fmt"

	"github.com/cheld/cicd-bot/pkg/config"
)

type Cli struct {
	config config.Configuration
}

func NewCli(config config.Configuration) *Cli {
	return &Cli{config}
}

func (cli *Cli) HandleStdin(args, stdin string) []config.DestinationCtx {
	triggers := []config.DestinationCtx{}
	for _, rule := range cli.config.Events.Cli.Stdin {
		event := config.EventCtx{args, nil}
		if rule.IsMatching(event) {
			fmt.Println(rule.IsMatching(event))
		}
		triggers = append(triggers, rule.DestinationCtx(event))
	}
	return triggers
}
