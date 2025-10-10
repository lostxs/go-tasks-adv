package request

import (
	"3-validation-api/pkg/response"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &payload, nil
}
