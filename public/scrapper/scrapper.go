// Package scrapper парсит урл и возвращает данные на превью
package scrapper

import (
	// Interfaces
	"MindAssistantBackend/interfaces/scrapper"
	// Helpers
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/helpers/response"
	// Packages
	"encoding/json"
	"net/http"
	// Libraries
	"github.com/badoux/goscraper"
)

// Parse парсит урл
func Parse(w http.ResponseWriter, r *http.Request) {
	// Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	request := new(iScrapper.Request)
	err := decoder.Decode(&request)
	errors.ErrorHandler(err, 500, w)
	defer r.Body.Close()

	if request.URL == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Парсинг ссылки
	scrap, err := goscraper.Scrape(request.URL, 5)
	if err != nil {
		errors.ErrorHandler(err, 500, w)
		return
	}

	// Формирование выходного значения
	result := &iScrapper.Response{
		Title:       scrap.Preview.Title,
		Description: scrap.Preview.Description,
		Image:       scrap.Preview.Images[0],
		Link:        scrap.Preview.Link,
	}

	// Формирование ответа от сервера
	response := response.Set(true, "", result)

	// Подготовка выходных данных
	output, err := json.Marshal(response)
	errors.ErrorHandler(err, 500, w)

	// Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}
