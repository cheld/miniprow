package webhook

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/cheld/cicd-bot/pkg/config"

	"github.com/cheld/cicd-bot/pkg/event"
	"github.com/cheld/cicd-bot/pkg/trigger"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	pathGithubWebhook = "/webhooks/github"
	pathHttpWebhook   = "/webhooks/http/"
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

	// Github
	http.HandleFunc(pathGithubWebhook, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Gihub event received")
		payload, err := githubWebhook.Parse(r, github.IssueCommentEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Printf("Github event not implemented.")
			} else {
				log.Printf("Error reading body: %s", err)
				http.Error(w, "can't read body", http.StatusBadRequest)
			}
			return
		}
		go dispatcher.Execute(handler.HandleGithub(payload))
	})

	// Http generic
	http.HandleFunc(pathHttpWebhook, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Http event received")
		path := strings.TrimPrefix(r.URL.Path, pathHttpWebhook)
		if path == "" {
			log.Printf("Call http webhook with url %s<event-type>", pathHttpWebhook)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		go dispatcher.Execute(handler.HandleHttp(body, path))
	})

	addr := fmt.Sprintf("%s:%d", opts.Bindaddr, opts.Port)
	fmt.Printf("Webhook listening on %s\n", addr)
	http.ListenAndServe(addr, nil)

}
