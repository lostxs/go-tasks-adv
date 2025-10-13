package resp

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, data any, statusCode int) {
	//устанавливаем ответ в виде json
	w.Header().Set("Content-Type", "application/json")
	//записыввем статус код ответа
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
