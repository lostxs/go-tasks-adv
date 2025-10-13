package req

import (
	"4-order-api/pkg/resp"
	"encoding/json"
	"io"
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		resp.Json(w, err.Error(), 402)
		return nil, err
	}
	return &body, nil
}
func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {

		return payload, err

	}
	return payload, nil
}
