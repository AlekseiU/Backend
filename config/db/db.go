// Package db хранит настройки БД
package db

import (
	// Interfaces
	"MindAssistantBackend/interfaces/config/db"
)

// Local содержит локальные доступы к БД
var Local = &iDb.Model{
	DbUser:     "urivsky",
	DbPassword: "123581321",
	DbName:     "mindassistant",
	DbSsl:      "disable",
}
