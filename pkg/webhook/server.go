package webhook

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cheld/cicd-bot/pkg/config"
	"github.com/cheld/cicd-bot/pkg/event"
	"github.com/cheld/cicd-bot/pkg/trigger"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	pathGithubWebhook = "/webhooks/github"
)

func Run() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	handler := event.NewHandler(cfg)
	dispatcher := trigger.NewDispatcher(cfg)

	githubWebhook, _ := github.New(github.Options.Secret("MySecret"))

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
	http.ListenAndServe(":3000", nil)
}
