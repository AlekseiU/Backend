// Package config хранит настройки сервиса
package config

import (
	// Config
	"MindAssistantBackend/config/db"
	"MindAssistantBackend/config/mode"
	"MindAssistantBackend/config/server"
	// Interfaces
	"MindAssistantBackend/interfaces/config/db"
	"MindAssistantBackend/interfaces/config/server"
)

// Mode содержит тип сервера
var Mode = mode.Type()

// Server устанавливает настройки сервера
func Server() *iServer.Model {
	return server.Config(*Mode)
}

// Db устанавливает настройки сервера
func Db() *iDb.Model {
	return db.Config(*Mode)
}
