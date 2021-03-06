package main

import (
	"net/http"

	calman "github.com/nelsonleduc/calmanbot/handlers"
	"github.com/nelsonleduc/correctionbot/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CalmanRespond",
		"POST",
		"/bots/botHook",
		calman.BotHook(handlers.HandleCalman),
	},
}
