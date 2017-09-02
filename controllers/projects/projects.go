// Projects содержит набор служебных функция для работы с проектами
package projects

import (
	// Config
	"MindAssistantBackend/config"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	// Interfaces
	"MindAssistantBackend/interfaces/projects"
	// Libraries
	"database/sql"
	"encoding/json"
	"net/http"
	// Packages
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Соединение с БД
var db = config.DbConnect()

// List отображает список проектов
func List(w http.ResponseWriter, r *http.Request) {
	// Подготовка запроса
	rows, err := db.Query("SELECT * FROM projects")
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

	// Подготовка запроса
	row := db.QueryRow("SELECT * FROM projects WHERE id = $1", id)
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
	row := db.QueryRow("INSERT INTO projects(name, pages) VALUES($1, $2) RETURNING id", project.Name, project.Pages)
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

	// Подготовка запроса
	update, err := db.Prepare("UPDATE projects SET name = $1, pages = $2 WHERE id = $3")
	errors.ErrorHandler(err, 500, w)

	// Выполнение запроса
	result, err := update.Exec(project.Name, project.Pages, project.ID)
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

	// Подготовка запроса на удаление проекта
	delete, err := db.Prepare("DELETE FROM projects WHERE id = $1")
	errors.ErrorHandler(err, 500, w)

	// Выполнение запросов
	result, err := delete.Exec(id)
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
