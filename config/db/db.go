// Package db хранит настройки БД
package db

import (
	// Interfaces
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/interfaces/config/db"
	// Libraries
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config возвращает настройки сервера
func Config(mode string) *iDb.Model {
	config, err := ioutil.ReadFile("./config/" + mode + "/db.json")
	if err != nil {
		errors.ErrorHandler(err, 500, nil)
		os.Exit(1)
	}

	var settings *iDb.Model
	json.Unmarshal(config, &settings)

	return settings
}
