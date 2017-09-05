// Package server хранит настройки сервера
package server

import (
	// Interfaces
	"MindAssistantBackend/interfaces/config/server"
)

// Local содержит локальные настройки сервера
var Local = &iServer.Model{
	Port:           "3000",
	AllowedMethods: []string{"GET", "DELETE", "POST", "PUT", "OPTIONS"},
}
