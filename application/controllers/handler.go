package controllers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data" swaggertype:"object"`
}
type Empty struct{}
type ResponseError struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error message"`
}

func JsonError(status int, w http.ResponseWriter, msg string) {
	var response ResponseError
	response.Success = false
	response.Message = msg

	w.Header().Set("Content-Type", "application/json")
	switch status {
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 404:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	r, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(r)
}

func JsonSuccess[T any](status int, w http.ResponseWriter, value T) {
	var response Response
	response.Success = true
	response.Message = ""

	w.Header().Set("Content-Type", "application/json")
	switch status {
	case 201:
		w.WriteHeader(http.StatusCreated)
	case 200:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonData, err := json.Marshal(value)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
	}
	response.Data = jsonData
	json, err := convertObjetcToJson(response)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
	}

	w.Write([]byte(json))
}

func convertObjetcToJson[T any](value T) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
