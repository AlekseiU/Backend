// Package iScrapper содержит интерфейс ответа парсера урла
package iScrapper

// Request структура запроса интерфейса
type Request struct {
	URL string `json:"url"`
}

// Response структура ответа интерфейса
type Response struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}
