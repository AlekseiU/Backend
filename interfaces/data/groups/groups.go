// Package iFieldGroup содержит интерфейсы для работы с группами полей
package iFieldGroup

import (
	// Interfaces
	"MindAssistantBackend/interfaces/data/fields"
)

type Model struct {
	ID     int64           `json:"id"`
	Name   string          `json:"name"`
	Order  int             `json:"order"`
	Data   int64           `json:"data"`
	Fields []*iField.Model `json:"fields"`
}
