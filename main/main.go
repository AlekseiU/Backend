// ToDo:
// Создание функционала профилей
package main

import (
	// Config
	"MindAssistantBackend/config"
	// Controllers
	"MindAssistantBackend/controllers/data"
	"MindAssistantBackend/controllers/data/fields"
	"MindAssistantBackend/controllers/data/groups"
	"MindAssistantBackend/controllers/projects"
	// Libraries
	"net/http"
	// Packages
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	// Методы для работы с проектами
	router.HandleFunc("/projects", projects.List).
		Methods("GET")

	router.HandleFunc("/projects/{id}", projects.Item).
		Methods("GET")

	router.HandleFunc("/projects", projects.Create).
		Methods("OPTIONS", "POST")

	router.HandleFunc("/projects/{id}", projects.Update).
		Methods("OPTIONS", "PUT")

	router.HandleFunc("/projects/{id}", projects.Delete).
		Methods("OPTIONS", "DELETE")

	// Методы для работы с Data объектами
	router.HandleFunc("/data/project/{id}", data.List).
		Methods("GET")

	router.HandleFunc("/data", data.Create).
		Methods("OPTIONS", "POST")

	router.HandleFunc("/data/{id}", data.Update).
		Methods("OPTIONS", "PUT")

	router.HandleFunc("/data/{id}", data.Delete).
		Methods("OPTIONS", "DELETE")

	// Методы для работы с группами полей Data объектов
	router.HandleFunc("/data/groups", groups.Create).
		Methods("OPTIONS", "POST")

	router.HandleFunc("/data/groups/{id}", groups.Delete).
		Methods("OPTIONS", "DELETE")

	// // Методы для работы с группами полей Data объектов
	router.HandleFunc("/data/fields", fields.Create).
		Methods("OPTIONS", "POST")

	router.HandleFunc("/data/fields/{id}", fields.Delete).
		Methods("OPTIONS", "DELETE")

	// Запуск сервера
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "DELETE", "POST", "PUT", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":"+config.Port, handler)
}
