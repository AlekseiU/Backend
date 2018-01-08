// Package routes содержит список поддерживаемых роутов
package routes

import (
	// Config
	"MindAssistantBackend/config"
	// Controllers
	"MindAssistantBackend/controllers/auth"
	"MindAssistantBackend/controllers/data"
	"MindAssistantBackend/controllers/data/fields"
	"MindAssistantBackend/controllers/data/groups"
	"MindAssistantBackend/controllers/projects"
	"MindAssistantBackend/controllers/user"
	// Public
	"MindAssistantBackend/public/scrapper"
	// Packages
	"net/http"
	// Libraries
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Handle формирует роуты
func Handle() http.Handler {
	router := mux.NewRouter()

	// Методы для работы с проектами
	router.Handle("/projects", auth.Middleware.Handler(http.HandlerFunc(projects.List))).
		Methods("GET")

	router.Handle("/projects/{id}", auth.Middleware.Handler(http.HandlerFunc(projects.Item))).
		Methods("GET")

	router.Handle("/projects", auth.Middleware.Handler(http.HandlerFunc(projects.Create))).
		Methods("OPTIONS", "POST")

	router.Handle("/projects/{id}", auth.Middleware.Handler(http.HandlerFunc(projects.Update))).
		Methods("OPTIONS", "PUT")

	router.Handle("/projects/{id}", auth.Middleware.Handler(http.HandlerFunc(projects.Delete))).
		Methods("OPTIONS", "DELETE")

	// Методы для работы с Data объектами
	router.Handle("/data/project/{id}", auth.Middleware.Handler(http.HandlerFunc(data.List))).
		Methods("GET")

	router.Handle("/data", auth.Middleware.Handler(http.HandlerFunc(data.Create))).
		Methods("OPTIONS", "POST")

	router.Handle("/data/{id}", auth.Middleware.Handler(http.HandlerFunc(data.Update))).
		Methods("OPTIONS", "PUT")

	router.Handle("/data/{id}", auth.Middleware.Handler(http.HandlerFunc(data.Delete))).
		Methods("OPTIONS", "DELETE")

	// Методы для работы с группами полей Data объектов
	router.Handle("/data/groups", auth.Middleware.Handler(http.HandlerFunc(groups.Create))).
		Methods("OPTIONS", "POST")

	router.Handle("/data/groups/{id}", auth.Middleware.Handler(http.HandlerFunc(groups.Delete))).
		Methods("OPTIONS", "DELETE")

	// Методы для работы с группами полей Data объектов
	router.Handle("/data/fields", auth.Middleware.Handler(http.HandlerFunc(fields.Create))).
		Methods("OPTIONS", "POST")

	router.Handle("/data/fields/{id}", auth.Middleware.Handler(http.HandlerFunc(fields.Delete))).
		Methods("OPTIONS", "DELETE")

	// Вспомогательные публичные методы
	router.HandleFunc("/scrap", scrapper.Parse).
		Methods("OPTIONS", "POST")

	// Методы для работы с пользователями
	router.HandleFunc("/user/registration", user.Registration).
		Methods("OPTIONS", "POST")

	router.HandleFunc("/user/login", user.Login).
		Methods("OPTIONS", "POST")

	// Список роутов
	handler := cors.New(cors.Options{
		AllowedOrigins:   config.Server().AllowedOrigins,
		AllowedMethods:   config.Server().AllowedMethods,
		AllowCredentials: config.Server().AllowCredentials,
		AllowedHeaders:   config.Server().AllowedHeaders,
	}).Handler(router)

	return handler
}
