package event

import (
	"github.com/cheld/miniprow/pkg/piper/config"
)

type Handler struct {
	config config.Configuration
	env    map[string]string
}

func NewHandler(config config.Configuration, env map[string]string) *Handler {
	return &Handler{
		config: config,
		env:    env,
	}
}
