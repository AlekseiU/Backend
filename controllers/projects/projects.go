// Package projects содержит набор служебных функция для работы с проектами
package projects

import (
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
	"net/http"
	// Libraries
	"github.com/gorilla/mux"
)

// Соединение с БД
var db = connect.Db()

// List отображает список проектов
func List(w http.ResponseWriter, r *http.Request) {
	// Подготовка запроса
	rows, err := dbProjects.List()
	errors.ErrorHandler(err, 500, w)
	defer rows.Close()

	// Сбор данных из БД в структуру
	projects := make([]*iProjects.Model, 0)
	for rows.Next() {
		project := new(iProjects.Model)

		err := rows.Scan(&project.ID, &project.Name, &project.Pages)
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
	err := row.Scan(&project.ID, &project.Name, &project.Pages)
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
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	project := new(iProjects.Model)
	err := decoder.Decode(&project)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if project.Name == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

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
