package webhook

import (
	"fmt"
	"net/http"

	"github.com/cheld/cicd-bot/pkg/config"

	"github.com/cheld/cicd-bot/pkg/event"
	"github.com/cheld/cicd-bot/pkg/trigger"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	pathGithubWebhook = "/webhooks/github"
)

type Options struct {
	Port     int
	Bindaddr string
	Secret   string
}

func Run(cfg config.Configuration, env map[string]string, opts Options) {

	handler := event.NewHandler(cfg, env)
	dispatcher := trigger.NewDispatcher(cfg)
	githubWebhook, _ := github.New(github.Options.Secret(opts.Secret))

	http.HandleFunc(pathGithubWebhook, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Gihub event received")
		payload, err := githubWebhook.Parse(r, github.IssueCommentEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		go dispatcher.Execute(handler.HandleGithub(payload))
	})

	addr := fmt.Sprintf("%s:%d", opts.Bindaddr, opts.Port)
	fmt.Printf("Webhook listening on %s", addr)
	http.ListenAndServe(addr, nil)

}
