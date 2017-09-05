// Package db хранит настройки БД
package db

import (
	// Interfaces
	"MindAssistantBackend/interfaces/config/db"
)

// Config содержит доступы к БД
var Config = &iDb.Model{
	DbUser:     "urivsky",
	DbPassword: "123581321",
	DbName:     "mindassistant",
	DbSsl:      "disable",
}
