// Package connect устанавливает соединение с БД
package connect

import (
	// Config
	"MindAssistantBackend/config"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	// Libraries
	"database/sql"
	// Регистрация драйвера Postgres
	_ "github.com/lib/pq"
)

// Db устанавливает соединение с БД
func Db() *sql.DB {
	config := config.Db()

	// Открываем соединение
	db, err := sql.Open("postgres", "user="+config.DbUser+" password="+config.DbPassword+" dbname="+config.DbName+" sslmode="+config.DbSsl)
	errors.ErrorHandler(err, 500, nil)

	// Отслеживаем состояние канала передачи данных
	errors.ErrorHandler(db.Ping(), 500, nil)

	return db
}
