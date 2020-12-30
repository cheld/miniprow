package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cheld/miniprow/pkg/piper/config"
	"github.com/cheld/miniprow/pkg/piper/event"
	"github.com/cheld/miniprow/pkg/piper/trigger"
	"github.com/golang/glog"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func NewPiper(piperCfg *[]byte, secret string) *Piper {
	cfg, err := config.Load(piperCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handler := event.NewHandler(cfg)
	dispatcher := trigger.NewDispatcher(cfg)
	githubWebhook, _ := github.New(github.Options.Secret(secret))

	piper := &Piper{
		mux: http.NewServeMux(),
	}

	piper.mux.Handle("/piper/github", handleGithub(handler, dispatcher, githubWebhook))
	piper.mux.Handle("/piper/http/", handleHTTP(handler, dispatcher))

	return piper
}

type Piper struct {
	mux *http.ServeMux
}

func (piper *Piper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	piper.mux.ServeHTTP(writer, request)
}

func handleGithub(handler *event.Handler, dispatcher trigger.Dispatcher, githubWebhook *github.Webhook) http.HandlerFunc {
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
		dispatcher.Execute(handler.HandleGithub(payload))
	}
}

func handleHTTP(handler *event.Handler, dispatcher trigger.Dispatcher) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Http event received")
		path := strings.TrimPrefix(req.URL.Path, "/webhooks/http/")
		if path == "" {
			glog.Errorf("Call http webhook with url %s<event-type>", "/webhooks/http/")
			http.Error(res, "can't read body", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			glog.Errorf("Error reading body: %v", err)
			http.Error(res, "can't read body", http.StatusBadRequest)
			return
		}
		go dispatcher.Execute(handler.HandleHttp(body, path))
	}
}
