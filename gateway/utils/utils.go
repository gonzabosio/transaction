package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, payload map[string]interface{}, statusCode int) {
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(payload)
}
