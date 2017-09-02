// Package config хранит настройки сервиса
package config

import (
	"MindAssistantBackend/helpers/errors"
	"database/sql"
)

// Port указывает на текущий порт бэкенда
const Port string = "3000"

// DbUser это пользователь БД
const DbUser string = "urivsky"

// DbPassword это пароль от БД
const DbPassword string = "123581321"

// DbName это имя БД
const DbName string = "mindassistant"

// DbSsl это активация режима ssl
const DbSsl string = "disable"

// DbConnect устанавливает соединение с БД
func DbConnect() *sql.DB {
	// Служебные переменные
	var err error
	var db *sql.DB

	// Открываем соединение
	db, err = sql.Open("postgres", "user="+DbUser+" password="+DbPassword+" dbname="+DbName+" sslmode="+DbSsl)
	errors.ErrorHandler(err, 500, nil)

	// Отслеживаем состояние канала передачи данных
	errors.ErrorHandler(db.Ping(), 500, nil)

	return db
}
