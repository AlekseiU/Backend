// Package dbUser содержит набор набор запросов для работы с пользователями
package dbUser

import (
	// Interfaces
	"MindAssistantBackend/interfaces/user"
	// Queries
	"MindAssistantBackend/db/connect"
	// Packages
	"database/sql"
)

// Соединение с БД
var db = dbConnect.Init()

// Get возращает пользователя по его email
func Get(email string) *sql.Row {
	return db.QueryRow("SELECT * FROM users WHERE email = $1", email)
}

// Create создает нового пользователя
func Create(user *iUser.Model) *sql.Row {
	return db.QueryRow("INSERT INTO users(email, password) VALUES($1, $2) RETURNING id", user.Email, user.Password)
}
