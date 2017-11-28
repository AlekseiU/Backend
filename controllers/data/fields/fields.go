// Package fields содержит набор служебных функций для работы с полями Data объектов
package fields

import (

	// Helpers
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	"strings"
	// Interfaces
	"MindAssistantBackend/interfaces/data/fields"
	"MindAssistantBackend/interfaces/data/groups"
	// Queries
	"MindAssistantBackend/db/connect"
	"MindAssistantBackend/db/data/fields"
	// Packages
	"encoding/json"
	"net/http"
	// Libraries
	"github.com/gorilla/mux"
)

// Соединение с БД
var db = dbConnect.Init()

// List выводит список полей
func List(w http.ResponseWriter, r *http.Request, group *iFieldGroup.Model) []*iField.Model {
	// Сбор и анализ входных данных
	if group.ID <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return nil
	}

	// Выполнение запроса
	rows, err := dbFields.List(group.ID)
	errors.ErrorHandler(err, 500, w)
	defer rows.Close()

	// Сбор данных из БД в структуру
	list := make([]*iField.Model, 0)
	for rows.Next() {
		field := new(iField.Db)

		err := rows.Scan(&field.ID, &field.Type, &field.Order, &field.Value, &field.Group, &field.Title)
		errors.ErrorHandler(err, 500, w)

		output := &iField.Model{
			ID:    field.ID,
			Type:  field.Type,
			Order: field.Order,
			Group: field.Group,
			Title: field.Title,
			Value: strings.Split(field.Value, ";"),
		}

		if output.Value == nil {
			output.Value = make([]string, 1)
		}

		list = append(list, output)
	}
	errors.ErrorHandler(rows.Err(), 500, w)

	return list
}

// Update обновляет поле
func Update(w http.ResponseWriter, r *http.Request, field *iField.Model) *iField.Db {
	if field.ID <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return nil
	}

	// Подготовка запроса
	output := &iField.Db{
		ID:    field.ID,
		Type:  field.Type,
		Order: field.Order,
		Group: field.Group,
		Title: field.Title,
	}

	if len(field.Value) > 1 {
		output.Value = strings.Join(field.Value, ";")
	} else {
		output.Value = field.Value[0]
	}

	// Выполнение запросов
	result, err := dbFields.Update(output)
	errors.ErrorHandler(err, 500, w)

	rows, err := result.RowsAffected()
	errors.ErrorHandler(err, 500, w)

	// Проверка на успешность
	if rows > 0 {
		// Отображение результата
		return output
	}

	return nil
}

// Create создает поле
func Create(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	field := new(iField.Model)
	err := decoder.Decode(&field)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if field.Group <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запроса
	row := dbFields.Create(field)
	err = row.Scan(&field.ID)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	if field.ID > 0 {
		// Формирование ответа от сервера
		response := response.Set(true, "", field)

		// Подготовка выходных данных
		output, err := json.Marshal(response)
		errors.ErrorHandler(err, 500, w)

		// Отображение результата
		w.Write(output)
	}
}

// Delete удаляет поле
func Delete(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Выполнение запросов
	result, err := dbFields.Delete(id)
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
