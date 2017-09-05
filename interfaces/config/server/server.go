// Package iServer содержит интерфейс настройки сервера
package iServer

// Model основная структура интерфейса
type Model struct {
	Port           string
	AllowedMethods []string
}
