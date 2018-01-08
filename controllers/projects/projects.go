// Package projects содержит набор служебных функция для работы с проектами
package projects

import (
	// Config
	"MindAssistantBackend/config"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	// Interfaces
	"MindAssistantBackend/interfaces/projects"
	// Queries
	"MindAssistantBackend/db/connect"
	"MindAssistantBackend/db/projects"
	// Packages
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	// Libraries
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// Соединение с БД
var db = dbConnect.Init()

// Токен пользователя
var tokenString string

// List отображает список проектов
func List(w http.ResponseWriter, r *http.Request) {
	// Получение токена из заголовков
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		tokenString = tokens[0]
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Проверка токета
	if tokenString == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Обработка токета
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.Server().Secret, nil
	})
	errors.ErrorHandler(err, 500, w)

	// Получение подписей из токена
	claims, _ := token.Claims.(jwt.MapClaims)

	// Подготовка запроса
	rows, err := dbProjects.List(claims["uid"].(float64))
	errors.ErrorHandler(err, 500, w)
	defer rows.Close()

	// Сбор данных из БД в структуру
	projects := make([]*iProjects.Model, 0)
	for rows.Next() {
		project := new(iProjects.Model)

		err := rows.Scan(&project.ID, &project.Name, &project.Pages, &project.Owner)
		errors.ErrorHandler(err, 500, w)

		projects = append(projects, project)
	}
	errors.ErrorHandler(rows.Err(), 500, w)

	// Формирование ответа от сервера
	response := response.Set(true, "", projects)

	// Подготовка выходных данных
	output, err := json.Marshal(response)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Item отображет проект по его id
func Item(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запроса
	row := dbProjects.Item(id)
	project := new(iProjects.Model)

	// Сбор данных из БД в структуру
	err := row.Scan(&project.ID, &project.Name, &project.Pages, &project.Owner)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}

	// Формирование ответа от сервера
	response := response.Set(true, "", project)

	// Подготовка выходных данных
	output, err := json.Marshal(response)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Create создает новый проект
func Create(w http.ResponseWriter, r *http.Request) {
	// Получение токена из заголовков
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		tokenString = tokens[0]
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Проверка токета
	if tokenString == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Обработка токета
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.Server().Secret, nil
	})
	errors.ErrorHandler(err, 500, w)

	// Получение подписей из токена
	claims, _ := token.Claims.(jwt.MapClaims)

	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	project := new(iProjects.Model)
	err = decoder.Decode(&project)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if project.Name == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Присвоение владельца проекту
	project.Owner = claims["uid"].(float64)

	// Выполнение запроса
	row := dbProjects.Create(project)
	err = row.Scan(&project.ID)
	errors.ErrorHandler(err, 500, w)

	if project.ID > 0 {
		// Формирование ответа от сервера
		response := response.Set(true, "", project)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(output)
	}
}

// Update обновляет данные проекта
func Update(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	project := new(iProjects.Model)
	err := decoder.Decode(&project)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if project.ID == 0 || project.Name == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запроса
	result, err := dbProjects.Update(project)
	errors.ErrorHandler(err, 500, w)

	// Проверка на успешность
	rows, err := result.RowsAffected()
	errors.ErrorHandler(err, 500, w)

	if rows > 0 {
		// Формирование ответа от сервера
		response := response.Set(true, "", project)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(output)
	}
}

// Delete удаляет проект по его id и все связанные с ним данные
func Delete(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запросов
	result, err := dbProjects.Delete(id)
	errors.ErrorHandler(err, 500, w)

	// Проверка на успешность
	rows, err := result.RowsAffected()
	errors.ErrorHandler(err, 500, w)

	if rows > 0 {
		// Формирование ответа от сервера
		response := response.Set(true, "", nil)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(output)
	}
}
