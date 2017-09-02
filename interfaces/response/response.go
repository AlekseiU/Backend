// Package iResponse содержит интерфейсы ответа клиенту
package iResponse

// Response - структура ответа
type Model struct {
	Result  bool        `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
