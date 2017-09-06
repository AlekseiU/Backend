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
	switch *Mode {
	case "local":
		return server.Local
	}

	return nil
}

// Db устанавливает настройки сервера
func Db() *iDb.Model {
	switch *Mode {
	case "local":
		return db.Local
	}

	return nil
}
