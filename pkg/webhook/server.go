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

func Run() {
	config := config.Load("config.yaml")

	handler := event.NewHandler(config)
	dispatcher := trigger.NewDispatcher(config)

	githubWebhook, _ := github.New(github.Options.Secret("MySecret"))

	http.HandleFunc(pathGithubWebhook, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Gihub event received")
		payload, err := githubWebhook.Parse(r, github.IssueCommentEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		dispatcher.Execute(handler.HandleGithub(payload))
	})
	http.ListenAndServe(":3000", nil)
}
