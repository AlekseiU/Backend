package main

import (
	// Config
	"MindAssistantBackend/config"
	// Routes
	"MindAssistantBackend/routes"
	// Libraries
	"net/http"
)

func main() {
	http.ListenAndServe(":"+config.Port, routes.Handle())
}
