package handlers

import (
	"regexp"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
)

func HandleCalman(message service.Message, service service.Service, cache cache.QueryCache) {

	if message.UserType() != "user" {
		return
	}

	bot, _ := models.FetchBot(message.GroupID())

	fix := firstCorrected(message.Text())
	if fix == "" {
		return
	}

	service.PostText(bot.Key, fix, -1, message)
}

// Private

var correctors []func(string) string

func init() {
	correctors = []func(string) string{
		CorrectImgur,
	}
}

func matchFunc(regex string) func(string) []string {
	return func(text string) []string {
		r, _ := regexp.Compile(regex)
		return r.FindStringSubmatch(text)
	}
}

func firstCorrected(text string) string {
	for _, f := range correctors {
		change := f(text)
		if change != "" {
			return change
		}
	}

	return ""
}
