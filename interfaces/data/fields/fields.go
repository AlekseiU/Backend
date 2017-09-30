// Package iField содержит интерфейсы для работы с полями
package iField

// Model основная структура интерфейса
type Model struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Order int    `json:"order"`
	Group int64  `json:"group"`
	Title string `json:"title"`
}
