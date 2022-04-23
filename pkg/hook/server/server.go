package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	config "github.com/cheld/miniprow/pkg/hook/model"
	_ "github.com/cheld/miniprow/pkg/hook/plugins-imports"
	"github.com/cheld/miniprow/pkg/hook/plugins/actions"
	"github.com/cheld/miniprow/pkg/hook/plugins/triggers"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func NewHandler(hookCfg *[]byte, secret string) *Hook {
	cfg, err := config.Load(hookCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Println(string(s))

	githubWebhook, _ := github.New(github.Options.Secret(secret))

	hook := &Hook{
		mux: http.NewServeMux(),
	}

	hook.mux.Handle("/hook/github", handleGithub(githubWebhook))
	hook.mux.Handle("/hook/http/", handleHTTP(cfg))

	return hook
}

type Hook struct {
	mux *http.ServeMux
}

func (piper *Hook) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	piper.mux.ServeHTTP(writer, request)
}

func handleGithub(githubWebhook *github.Webhook) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Gihub event received")
		payload, err := githubWebhook.Parse(req, github.IssueCommentEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				glog.Infof("Github event not implemented.")
			} else {
				glog.Errorf("Error reading body: %s", err)
				http.Error(res, "can't read body", http.StatusBadRequest)
			}
			return
		}
		event := config.Event{}
		event.Data = payload
		tenant := config.Tenant{}
		triggeredRules := triggers.Handle(&event, tenant)
		actions.Handle(triggeredRules, &event, tenant)
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
		triggeredRules := triggers.Handle(&event, tenant)
		actions.Handle(triggeredRules, &event, tenant)

		fmt.Printf("%s", event.Trail())
	}
}
