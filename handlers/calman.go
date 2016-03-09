package handlers

import (
	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/service"
)

func HandleCalman(message service.Message, service service.Service, cache cache.QueryCache) {

	if message.UserType() != "user" {
		return
	}

	// bot, _ := models.FetchBot(message.GroupID())
}
