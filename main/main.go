package main

import (
	// Config
	"MindAssistantBackend/config"
	"os"
	// Routes
	"MindAssistantBackend/routes"
	// Packages
	"net/http"
	// Libraries
	"github.com/gorilla/handlers"
)

func main() {
	// Запуск сервера
	http.ListenAndServe(":"+config.Server().Port, handlers.LoggingHandler(os.Stdout, routes.Handle()))
}
