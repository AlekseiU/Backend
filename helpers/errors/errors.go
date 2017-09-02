// Package errors обрабатывает ошибки
package errors

import (
	"log"
	"net/http"
)

// ErrorHandler проверяет ошибки
func ErrorHandler(err error, code int, w http.ResponseWriter) {
	if err != nil {
		log.Print(err) // <= Режим отладки
		if w != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}
