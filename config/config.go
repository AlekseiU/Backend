// Package config хранит настройки сервиса
package config

import (
	"MindAssistantBackend/helpers/errors"
	"database/sql"
	// Interfaces
	"MindAssistantBackend/interfaces/config"
	// Регистрация драйвера Postgres
	_ "github.com/lib/pq"
)

// Server устанавливает настройки сервера
func Server(mode *string) *iConfig.Server {
	switch *mode {
	case "local":
		config := &iConfig.Server{
			Port:           "3000",
			AllowedMethods: []string{"GET", "DELETE", "POST", "PUT", "OPTIONS"},
		}

		return config
	}

	return nil
}

// Db устанавливает настройки БД
func Db() *iConfig.Db {
	config := &iConfig.Db{
		DbUser:     "urivsky",
		DbPassword: "123581321",
		DbName:     "mindassistant",
		DbSsl:      "disable",
	}

	return config
}

// DbConnect устанавливает соединение с БД
func DbConnect() *sql.DB {
	// Служебные переменные
	var err error
	var db *sql.DB

	// Открываем соединение
	db, err = sql.Open("postgres", "user="+Db().DbUser+" password="+Db().DbPassword+" dbname="+Db().DbName+" sslmode="+Db().DbSsl)
	errors.ErrorHandler(err, 500, nil)

	// Отслеживаем состояние канала передачи данных
	errors.ErrorHandler(db.Ping(), 500, nil)

	return db
}
