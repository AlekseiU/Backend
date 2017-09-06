// Package dbProjects содержит набор набор запросов для работы с проектами
package dbProjects

import (
	// Interfaces
	"MindAssistantBackend/interfaces/projects"
	// Queries
	"MindAssistantBackend/db/connect"
	// Packages
	"database/sql"
)

// Соединение с БД
var db = connect.Db()

// List отображает список проектов
func List() (*sql.Rows, error) {
	return db.Query("SELECT * FROM projects")
}

// Item отображет проект по его id
func Item(id string) *sql.Row {
	return db.QueryRow("SELECT * FROM projects WHERE id = $1", id)
}

// Create создает новый проект
func Create(project *iProjects.Model) *sql.Row {
	return db.QueryRow("INSERT INTO projects(name, pages) VALUES($1, $2) RETURNING id", project.Name, project.Pages)
}

// Update обновляет данные проекта
func Update(project *iProjects.Model) (sql.Result, error) {
	// Подготовка запроса
	update, err := db.Prepare("UPDATE projects SET name = $1, pages = $2 WHERE id = $3")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return update.Exec(project.Name, project.Pages, project.ID)
}

// Delete удаляет проект по его id и все связанные с ним данные
func Delete(id string) (sql.Result, error) {
	// Подготовка запроса
	delete, err := db.Prepare("DELETE FROM projects WHERE id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return delete.Exec(id)
}
