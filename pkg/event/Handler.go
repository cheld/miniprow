package event

import (
	"github.com/cheld/cicd-bot/pkg/config"
)

type Handler struct {
	config config.Configuration
}

func NewHandler(config config.Configuration) *Handler {
	return &Handler{
		config: config}
}
