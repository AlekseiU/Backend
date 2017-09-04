// Package iResponse содержит интерфейсы ответа клиенту
package iResponse

// Model основная структура интерфейса
type Model struct {
	Result  bool        `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
