// Package iProjects содержит интерфейсы для работы с проектами
package iProjects

// Model основная структура интерфейса
type Model struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Pages int     `json:"pages"`
	Owner float64 `json:"owner"`
}
