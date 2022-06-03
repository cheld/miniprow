package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	config "github.com/cheld/miniprow/pkg/hook/model"
	"github.com/cheld/miniprow/pkg/hook/rules"
	_ "github.com/cheld/miniprow/pkg/hook/rulesimports"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/sirupsen/logrus"
)

func NewHandler(hookCfg *[]byte, secret string) *Hook {
	cfg, err := config.Load(hookCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Println(string(s))

	//githubWebhook, _ := github.New(github.Options.Secret(""))
	githubWebhook, _ := github.New(github.Options.Secret("asdf"))

	hook := &Hook{
		mux: http.NewServeMux(),
	}

	hook.mux.Handle("/hook/github", handleGithub(githubWebhook, cfg))
	hook.mux.Handle("/hook/http/", handleHTTP(cfg))

	return hook
}

type Hook struct {
	mux *http.ServeMux
}

func (piper *Hook) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	piper.mux.ServeHTTP(writer, request)
}

func handleGithub(githubWebhook *github.Webhook, cfg config.Configuration) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logrus.Infof("Github event received")
		tenant := config.Tenant{}
		tenant.Config = cfg
		event := config.Event{}
		payload, err := githubWebhook.Parse(req, github.IssueCommentEvent, github.InstallationEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				logrus.Infof("Github event not implemented: %s", err)
			} else {
				logrus.Errorf("Error reading body: %s", err)
				//logrus.Error(res, "can't read body", http.StatusBadRequest)
			}
			return
		}
		event.Data = payload
		//fmt.Println(event.Data.(github.InstallationPayload).Installation.Account.ReposURL)
		//fmt.Println(event.Data.(github.InstallationPayload).Installation.RepositorySelection)
		//fmt.Println(event.Data.(github.InstallationPayload).Repositories[0].FullName)
		//fmt.Println(event.Data.(github.InstallationPayload).Installation.Account.URL)
		//fmt.Println(event.Data.(github.InstallationPayload).Installation.Account.Login)
		//fmt.Println(event.Data.(github.IssueCommentPayload).Repository.FullName)
		//fmt.Println(event.Data.(github.IssueCommentPayload).Repository.Owner.SiteAdmin)
		//fmt.Println(event.Data.(github.IssueCommentPayload).Repository.Name)
		event.Type = "github_comment"

		listeners := rules.NewRuleBasedListeners(cfg.Rules)
		for _, l := range listeners {
			l.Handle(event)
		}

		fmt.Printf("%s", event.Trail())
	}
}

func handleHTTP(cfg config.Configuration) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logrus.Infof("Http event received")
		tenant := config.Tenant{}
		tenant.Config = cfg

		//parse event
		event := config.Event{}
		event.Type = "http"
		path := strings.TrimPrefix(req.URL.Path, "/hook/http/")
		if path != "" {
			event.Type = "http_" + path
		}
		event.Log("Event type: %v", event.Type)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			event.Log("Error reading body: %v", err)
			http.Error(res, "can't read body", http.StatusBadRequest)
			return
		}
		event.Data = body

		// execute

		listeners := rules.NewRuleBasedListeners(cfg.Rules)
		for _, l := range listeners {
			l.Handle(event)
		}

		fmt.Printf("%s", event.Trail())
	}
}
