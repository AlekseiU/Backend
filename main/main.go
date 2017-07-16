// ToDo:
// Создание функционала профилей
package main

import (
	"MindAssistantBackend/data"
	"MindAssistantBackend/projects"
	"net/http"
)

func main() {
	// Методы для работы с проектами
	http.HandleFunc("/projects", projects.List)
	http.HandleFunc("/projects/show", projects.Item)
	http.HandleFunc("/projects/create", projects.Create)
	http.HandleFunc("/projects/update", projects.Update)
	http.HandleFunc("/projects/delete", projects.Delete)

	// Методы для работы с Data объектами
	http.HandleFunc("/data", data.List)
	http.HandleFunc("/data/project", data.ListByProject)
	http.HandleFunc("/data/show", data.Item)
	http.HandleFunc("/data/create", data.Create)
	http.HandleFunc("/data/update", data.Update)
	http.HandleFunc("/data/delete", data.Delete)

	// Запуск сервера
	http.ListenAndServe(":3000", nil)
}
