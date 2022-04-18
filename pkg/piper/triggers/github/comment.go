package github

import "github.com/cheld/miniprow/pkg/piper/triggers"

const (
	triggerName = "github_comment"
)

func init() {
	triggers.RegisterGithubHandler(triggerName, handleEvent)
}
