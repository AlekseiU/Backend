package main

import (
	// Config
	"MindAssistantBackend/config"
	// Routes
	"MindAssistantBackend/routes"
	// Libraries
	"flag"
	"net/http"
)

func main() {
	// Тип запускаемого сервера
	mode := flag.String("mode", "", "")
	flag.Parse()

	// Запуск сервера
	http.ListenAndServe(":"+config.Server(mode).Port, routes.Handle(mode))
}
