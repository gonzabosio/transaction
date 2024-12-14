package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func WriteJSON(rw http.ResponseWriter, payload interface{}, statusCode int) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(payload)
}

func ValidateStruct(body interface{}) error {
	if err := validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}
	return nil
}
