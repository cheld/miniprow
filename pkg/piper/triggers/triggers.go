package triggers

import "github.com/cheld/miniprow/pkg/piper/config"

var (
	githubHandlers = map[string]GithubHandler{}
)

// GithubHandler defines the function contract for a github handler.
type GithubHandler func(config.Configuration, config.Event) config.Rule

func RegisterGithubHandler(name string, fn GithubHandler) {
	githubHandlers[name] = fn
}

func Handle(event config.Event) {

}
