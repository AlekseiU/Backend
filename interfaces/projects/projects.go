// Package iProjects содержит интерфейсы для работы с проектами
package iProjects

// Project - структура проекта
type Model struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Pages int    `json:"pages"`
}
