// Package server хранит настройки сервера
package server

import (
	// Interfaces
	"MindAssistantBackend/helpers/errors"
	"MindAssistantBackend/interfaces/config/server"
	// Libraries
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config возвращает настройки сервера
func Config(mode string) *iServer.Model {
	raw, err := ioutil.ReadFile("./config/" + mode + "/server.json")
	if err != nil {
		errors.ErrorHandler(err, 500, nil)
		os.Exit(1)
	}

	var settings *iServer.Model
	json.Unmarshal(raw, &settings)

	return settings
}
