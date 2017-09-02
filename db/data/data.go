// dbData содержит набор набор запросов для работы с Data
package dbData

import (
	// Config
	"MindAssistantBackend/config"
	// Interfaces
	"MindAssistantBackend/interfaces/data"
	// Packages
	"database/sql"

	_ "github.com/lib/pq"
)

// Соединение с БД
var db = config.DbConnect()

// List отображает список проектов
func List(id string) (*sql.Rows, error) {
	return db.Query("SELECT * FROM data WHERE project = $1", id)
}

// Create создает новый проект
func Create(data *iData.Model, coordinates []byte) *sql.Row {
	return db.QueryRow("INSERT INTO data(name, project, parent, coordinates) VALUES($1, $2, $3, $4) RETURNING id", data.Name, data.Project, data.Parent, coordinates)
}

// Update обновляет данные проекта
func Update(data *iData.Model, coordinates []byte) (sql.Result, error) {
	// Подготовка запроса
	update, err := db.Prepare("UPDATE data SET name = $2, project = $3, parent = $4, coordinates = $5 WHERE id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return update.Exec(data.ID, data.Name, data.Project, data.Parent, coordinates)
}

// Delete удаляет проект по его id и все связанные с ним данные
func Delete(id string) (sql.Result, error) {
	// Подготовка запроса
	delete, err := db.Prepare("DELETE FROM data WHERE id = $1")
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	return delete.Exec(id)
}
