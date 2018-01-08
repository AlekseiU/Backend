// Package iUser содержит интерфейс пользователя
package iUser

// Model основная структура интерфейса
type Model struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
	Token    string `json:"token"`
}

// Db структура интерфейса в БД
type Db struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

// JSON структура интерфейса для ответа клиенту
type JSON struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
