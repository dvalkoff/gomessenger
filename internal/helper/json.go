package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorDto struct {
    Code int `json:"code"`
    ErrorInfo string `json:"errorInfo"`
}

func EncodeError(w http.ResponseWriter, r *http.Request, status int, err error) error {
    errorDto := ErrorDto{
        Code: status,
        ErrorInfo: err.Error(),
    }
    return Encode(w, r, status, errorDto)
}

func EncodeNoBody(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        return fmt.Errorf("encode json: %w", err)
    }
    return nil
}

func Decode[T any](r *http.Request) (T, error) {
    var v T
    if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
        return v, fmt.Errorf("decode json: %w", err)
    }
    return v, nil
}
