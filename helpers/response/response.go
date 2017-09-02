// Package response формирует ответ сервера
package response

import (
	// Interfaces
	"MindAssistantBackend/interfaces/response"
)

// Set формирует ответ
func Set(result bool, message string, data interface{}) iResponse.Model {
	return iResponse.Model{
		Result:  result,
		Message: message,
		Data:    data,
	}
}
