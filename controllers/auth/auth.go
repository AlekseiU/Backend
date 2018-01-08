// Package auth содержит набор служебных функция для работы с токенами
package auth

import (
	// Config
	"MindAssistantBackend/config"
	// Libraries
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// Конфигурация сервера
var serverConfig = config.Server()

// Проверка на авторизационный заголовок
var Middleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return serverConfig.Secret, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
