// Package iData содержит интерфейсы для работы с Data объектами
package iData

import (
	// Interfaces

	"MindAssistantBackend/interfaces/data/groups"
)

// Model основная структура интерфейса
type Model struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Project     int                  `json:"project"`
	Parent      int                  `json:"parent"`
	Coordinates map[string]float64   `json:"coordinates"`
	Content     []*iFieldGroup.Model `json:"content"`
}

// Db структура для работы с БД
type Db struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Project     int    `json:"project"`
	Parent      int    `json:"parent"`
	Coordinates []byte `json:"coordinates"`
}
