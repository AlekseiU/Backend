// Package config хранит настройки сервиса
package config

import (
	// Config
	"MindAssistantBackend/config/db"
	"MindAssistantBackend/config/server"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	// Interfaces
	"MindAssistantBackend/interfaces/config/server"
	// Libraries
	"database/sql"
	// Регистрация драйвера Postgres
	_ "github.com/lib/pq"
)

// Server устанавливает настройки сервера
func Server(mode *string) *iServer.Model {
	switch *mode {
	case "local":
		return server.Local
	}

	return nil
}

// DbConnect устанавливает соединение с БД
func DbConnect() *sql.DB {
	// Открываем соединение
	output, err := sql.Open("postgres", "user="+db.Config.DbUser+" password="+db.Config.DbPassword+" dbname="+db.Config.DbName+" sslmode="+db.Config.DbSsl)
	errors.ErrorHandler(err, 500, nil)

	// Отслеживаем состояние канала передачи данных
	errors.ErrorHandler(output.Ping(), 500, nil)

	return output
}
