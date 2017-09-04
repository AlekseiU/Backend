// Package dbFieldGroups содержит набор запросов для работы с группами полей Data объектов
package dbFieldGroups

import (
	// Config
	"MindAssistantBackend/config"
	// Interfaces
	"MindAssistantBackend/interfaces/data/groups"
	// Packages
	"database/sql"
)

// Соединение с БД
var db = config.DbConnect()

// List отображает список проектов
func List(id int64) (*sql.Rows, error) {
	return db.Query("SELECT * FROM field_groups WHERE data = $1", id)
}

// Create создает новый проект
func Create(group *iFieldGroup.Model) *sql.Row {
	return db.QueryRow("INSERT INTO field_groups (name, ordr, data) VALUES ($1, $2, $3) RETURNING id", group.Name, group.Order, group.Data)
}

// Update обновляет данные проекта
func Update(group *iFieldGroup.Model) (sql.Result, error) {
	// Подготовка запроса
	update, err := db.Prepare("UPDATE field_groups set name = $2, ordr = $3, data = $4 where id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return update.Exec(group.ID, group.Name, group.Order, group.Data)
}

// Delete удаляет проект по его id и все связанные с ним данные
func Delete(id string) (sql.Result, error) {
	// Подготовка запроса
	delete, err := db.Prepare("DELETE FROM field_groups WHERE id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return delete.Exec(id)
}
