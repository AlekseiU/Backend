// Package user содержит набор служебных функция для работы с пользователями
package user

import (
	// Helpers
	"MindAssistantBackend/config"
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	// Interfaces
	"MindAssistantBackend/interfaces/user"
	// Queries
	"MindAssistantBackend/db/connect"
	"MindAssistantBackend/db/user"
	// Packages
	"database/sql"
	"encoding/json"
	"net/http"
	// Libraries
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Соединение с БД
var db = dbConnect.Init()

// Registration регистрирует пользователя
func Registration(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	user := new(iUser.Model)
	err := decoder.Decode(&user)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if user.Email == "" || user.Password == nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Генерируем hash на основе пароля
	user.Password, err = bcrypt.GenerateFromPassword(user.Password, 10)
	errors.ErrorHandler(err, 500, w)

	// Выполнение запроса
	row := dbUser.Create(user)
	err = row.Scan(&user.ID)
	errors.ErrorHandler(err, 500, w)

	if user.ID > 0 {
		// Создаем токен
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"uid": user.ID,
		})

		// Подписываем токен нашим секретным ключем
		accessToken, err := token.SignedString(config.Server().Secret)
		errors.ErrorHandler(err, 500, w)

		// Формирование ответа от сервера
		jsonUser := &iUser.JSON{
			Email: user.Email,
			Token: accessToken,
		}
		result := response.Set(true, "", jsonUser)

		// Подготовка выходных данных
		output, err := json.Marshal(result)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(output)
	}
}

// Login авторизует пользователя
func Login(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	user := new(iUser.Model)
	err := decoder.Decode(&user)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if user.Email == "" || user.Password == nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запроса
	row := dbUser.Get(user.Email)
	dbUser := new(iUser.Db)

	// Сбор данных из БД в структуру
	err = row.Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}

	// Формирование ответа от сервера
	jsonUser := &iUser.JSON{
		Email: dbUser.Email,
	}
	result := response.Set(true, "", jsonUser)

	err = bcrypt.CompareHashAndPassword(dbUser.Password, user.Password)
	if err != nil {
		result = response.Set(false, "Password incorrect", nil)
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": dbUser.ID,
	})

	// Подписываем токен нашим секретным ключем
	jsonUser.Token, err = token.SignedString(config.Server().Secret)
	errors.ErrorHandler(err, 500, w)

	// Подготовка выходных данных
	output, err := json.Marshal(result)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}
