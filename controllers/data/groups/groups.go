// Package groups содержит набор служебных функций для работы с группами полей Data объектов
package groups

import (

	// Config
	"MindAssistantBackend/config"
	"encoding/json"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	// Interfaces
	"MindAssistantBackend/interfaces/data"
	"MindAssistantBackend/interfaces/data/groups"
	// Queries
	"MindAssistantBackend/db/data/groups"
	// Controllers
	"MindAssistantBackend/controllers/data/fields"
	// Libraries
	"net/http"
	// Packages
	"github.com/gorilla/mux"
)

// Соединение с БД
var db = config.DbConnect()

// List выводит группы полей
func List(w http.ResponseWriter, r *http.Request, data *iData.Db) []*iFieldGroup.Model {
	// Сбор и анализ входных данных
	if data.ID <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return nil
	}

	// Выполнение запроса
	rows, err := dbFieldGroups.List(data.ID)
	errors.ErrorHandler(err, 500, w)
	defer rows.Close()

	// Сбор данных из БД в структуру
	list := make([]*iFieldGroup.Model, 0)
	for rows.Next() {
		group := new(iFieldGroup.Model)

		err := rows.Scan(&group.ID, &group.Name, &group.Order, &group.Data)
		errors.ErrorHandler(err, 500, w)

		// Сбор данных из таблицы fields связанных с Data объектом
		group.Fields = fields.List(w, r, group)

		list = append(list, group)
	}
	errors.ErrorHandler(rows.Err(), 500, w)

	return list
}

// Update обновляет группу полей
func Update(w http.ResponseWriter, r *http.Request, group *iFieldGroup.Model) *iFieldGroup.Model {
	if group.ID <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return nil
	}

	// Выполнение запросов
	result, err := dbFieldGroups.Update(group)
	errors.ErrorHandler(err, 500, w)

	rows, err := result.RowsAffected()
	errors.ErrorHandler(err, 500, w)

	// Обработка контента Data объекта
	if group.Fields != nil {
		// Обновление групп полей
		for f := 0; f < len(group.Fields); f++ {
			fields.Update(w, r, group.Fields[f])
		}
	}

	// Проверка на успешность
	if rows > 0 {
		// Отображение результата
		return group
	}

	return nil
}

// Create создает группу полей
func Create(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	group := new(iFieldGroup.Model)
	err := decoder.Decode(&group)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if group.Data <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запроса
	row := dbFieldGroups.Create(group)
	err = row.Scan(&group.ID)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	if group.ID > 0 {
		// Формирование ответа от сервера
		response := response.Set(true, "", group)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Write(output)
	}
}

// Delete удаляет группу полей
func Delete(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запросов
	result, err := dbFieldGroups.Delete(id)
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
		w.Write(output)
	}
}
