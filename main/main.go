package main

import (
	// Config
	"MindAssistantBackend/config"
	// Routes
	"MindAssistantBackend/routes"
	// Packages
	"net/http"
)

func main() {
	// Запуск сервера
	http.ListenAndServe(":"+config.Server().Port, routes.Handle())
}
