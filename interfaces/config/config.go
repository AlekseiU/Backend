// Package iConfig содержит интерфейс настройки сервера
package iConfig

// Server структура настроек сервера
type Server struct {
	Port           string
	AllowedMethods []string
}

// Db структура настроек БД
type Db struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbSsl      string
}
