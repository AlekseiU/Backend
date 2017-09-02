// Package data содержит набор служебных функций для работы с Data объектами
package data

import (
	// Config
	"MindAssistantBackend/config"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	// Interfaces
	"MindAssistantBackend/interfaces/data"
	"MindAssistantBackend/interfaces/data/groups"
	// Controllers
	"MindAssistantBackend/controllers/data/groups"
	// Libraries
	"encoding/json"
	"net/http"
	// Packages
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Соединение с БД
var db = config.DbConnect()

// List отображает список всех Data объектов
func List(w http.ResponseWriter, r *http.Request) {
	// Подготовка запроса
	dataRows, err := db.Query("SELECT * FROM data")
	errors.ErrorHandler(err, 500, w)
	defer dataRows.Close()

	// Сбор данных из БД в структуру
	dataList := make([]*iData.Model, 0)
	for dataRows.Next() {
		data := new(iData.Db)

		err := dataRows.Scan(&data.ID, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
		errors.ErrorHandler(err, 500, w)

		// Обработка координат
		var coordinates map[string]float64
		json.Unmarshal([]byte(data.Coordinates), &coordinates)

		// Сбор данных из таблицы field_group связанных с Data объектом
		content := groups.List(w, r, data)

		// Трансформация Data в новый объект
		dataResult := &iData.Model{
			ID:          data.ID,
			Name:        data.Name,
			Project:     data.Project,
			Parent:      data.Parent,
			Coordinates: coordinates,
			Content:     content,
		}

		dataList = append(dataList, dataResult)
	}
	errors.ErrorHandler(dataRows.Err(), 500, w)

	// Формирование ответа от сервера
	response := response.Set(true, "", dataList)

	// Подготовка выходных данных
	output, err := json.Marshal(response)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// ListByProject отображает список Data объектов по id проекта
func ListByProject(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	params := mux.Vars(r)
	project := params["id"]
	if project == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Подготовка запроса
	dataRows, err := db.Query("SELECT * FROM data WHERE project = $1", project)
	errors.ErrorHandler(err, 500, w)
	defer dataRows.Close()

	// Сбор данных из БД в структуру
	dataList := make([]*iData.Model, 0)
	for dataRows.Next() {
		data := new(iData.Db)

		err := dataRows.Scan(&data.ID, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
		errors.ErrorHandler(err, 500, w)

		// Обработка координат
		var coordinates map[string]float64
		json.Unmarshal([]byte(data.Coordinates), &coordinates)

		// Сбор данных из таблицы field_group связанных с Data объектом
		content := groups.List(w, r, data)

		// Трансформация Data в новый объект
		dataResult := &iData.Model{
			ID:          data.ID,
			Name:        data.Name,
			Project:     data.Project,
			Parent:      data.Parent,
			Coordinates: coordinates,
			Content:     content,
		}

		dataList = append(dataList, dataResult)
	}
	errors.ErrorHandler(dataRows.Err(), 500, w)

	// Формирование ответа от сервера
	response := response.Set(true, "", dataList)

	// Подготовка выходных данных
	output, err := json.Marshal(response)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Create создает новый Data объект
func Create(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	data := new(iData.Model)
	err := decoder.Decode(&data)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if data.Name == "" || data.Project <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Обработка координат
	coordinates, err := json.Marshal(data.Coordinates)
	errors.ErrorHandler(err, 500, w)

	// Выполнение запроса
	row := db.QueryRow("INSERT INTO data(name, project, parent, coordinates) VALUES($1, $2, $3, $4) RETURNING id", data.Name, data.Project, data.Parent, coordinates)
	err = row.Scan(&data.ID)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	if data.ID > 0 {
		// Обработка полей объекта
		content := make([]*iFieldGroup.Model, 0)

		// Трансформация Data в новый объект
		DataResult := &iData.Model{
			ID:          data.ID,
			Name:        data.Name,
			Project:     data.Project,
			Parent:      data.Parent,
			Coordinates: data.Coordinates,
			Content:     content,
		}

		// Формирование ответа от сервера
		response := response.Set(true, "", DataResult)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Write(output)
	}
}

// Update обновляет Data объект
func Update(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	data := new(iData.Model)
	err := decoder.Decode(&data)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if data.ID == 0 || data.Project <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Обработка координат
	if data.Coordinates == nil {
		data.Coordinates = map[string]float64{"x": 0, "y": 0}
	}
	coordinates, _ := json.Marshal(data.Coordinates)

	// Подготовка запроса
	update, err := db.Prepare("UPDATE data SET name = $2, project = $3, parent = $4, coordinates = $5 WHERE id = $1")
	errors.ErrorHandler(err, 500, w)

	// Выполнение запроса
	result, err := update.Exec(data.ID, data.Name, data.Project, data.Parent, coordinates)
	errors.ErrorHandler(err, 500, w)

	// Проверка на успешность
	rows, err := result.RowsAffected()
	errors.ErrorHandler(err, 500, w)

	// Обработка контента Data объекта
	if data.Content != nil {
		// Обновление групп полей
		for g := 0; g < len(data.Content); g++ {
			groups.Update(w, r, data.Content[g])
		}
	}

	// Проверка на успешность
	if rows > 0 {
		// Формирование ответа от сервера
		response := response.Set(true, "", data)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Write(output)
	}
}

// Delete удаляет Data объект по его id и все связанные с ним данные
func Delete(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Подготовка запроса на удаление Data объекта
	deleteData, err := db.Prepare("DELETE FROM data WHERE id = $1")
	errors.ErrorHandler(err, 500, w)

	// Выполнение запросов
	deleteData.Exec(id)

	// Формирование ответа от сервера
	response := response.Set(true, "", nil)

	// Подготовка выходных данных
	output, err := json.Marshal(response)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Write(output)
}
