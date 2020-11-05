package main

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

func main() {

	config := config.Load("config.yaml")
	args := "test"
	stdin := ""
	handler := event.NewHandler(config)
	dispatcher := trigger.NewDispatcher(config)
	eventData := handler.HandleCli(args, stdin)
	fmt.Println(eventData)
	trigger.NewDispatcher(config).Execute(eventData)

	//event := source.NewGithub(config).ParseInput(payload)
	// destination.New(config).Execute(event)

	githubWebhook, _ := github.New(github.Options.Secret("MySecret"))

	http.HandleFunc(pathGithubWebhook, func(w http.ResponseWriter, r *http.Request) {
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
