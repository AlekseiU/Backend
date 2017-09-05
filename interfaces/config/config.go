// Package iConfig содержит интерфейс настройки сервера
package iConfig

import (
	// Interfaces
	"MindAssistantBackend/interfaces/config/db"
	"MindAssistantBackend/interfaces/config/server"
)

// Model основная структура интерфейса
type Model struct {
	Server *iServer.Model
	Db     *iDb.Model
}
