// Package dbFields содержит набор набор запросов для работы с полями Data объектов
package dbFields

import (
	// Interfaces
	"MindAssistantBackend/interfaces/data/fields"
	// Queries
	"MindAssistantBackend/db/connect"
	// Packages
	"database/sql"
	// Libraries
	"github.com/lib/pq"
)

// Соединение с БД
var db = dbConnect.Init()

// List отображает список проектов
func List(id int64) (*sql.Rows, error) {
	return db.Query("SELECT * FROM fields WHERE group_id = $1", id)
}

// Create создает новый проект
func Create(field *iField.Model) *sql.Row {
	return db.QueryRow("INSERT INTO fields (type, ordr, value, group_id, title) VALUES ($1, $2, $3, $4, $5) RETURNING id", field.Type, field.Order, pq.Array(field.Value), field.Group, field.Title)
}

// Update обновляет данные проекта
func Update(field *iField.Db) (sql.Result, error) {
	// Подготовка запроса
	update, err := db.Prepare("UPDATE fields set type = $2, ordr = $3, value = $4, group_id = $5, title = $6 where id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return update.Exec(field.ID, field.Type, field.Order, field.Value, field.Group, field.Title)
}

// Delete удаляет проект по его id и все связанные с ним данные
func Delete(id string) (sql.Result, error) {
	// Подготовка запроса
	delete, err := db.Prepare("DELETE FROM fields WHERE id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return delete.Exec(id)
}
