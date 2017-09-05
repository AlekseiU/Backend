// Package iDb содержит интерфейс настройки БД
package iDb

// Model основная структура интерфейса
type Model struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbSsl      string
}
