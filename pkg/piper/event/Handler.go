package event

import (
	"github.com/cheld/miniprow/pkg/piper/config"
)

type Handler struct {
	config config.Configuration
}

func NewHandler(config config.Configuration) *Handler {
	return &Handler{
		config: config,
	}
}
