// Package iServer содержит интерфейс настройки сервера
package iServer

// Model основная структура интерфейса
type Model struct {
	Port             string
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	Secret           []byte
}
