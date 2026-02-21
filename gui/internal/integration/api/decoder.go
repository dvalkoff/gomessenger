package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorDto struct {
	Code      int    `json:"code"`
	ErrorInfo string `json:"errorInfo"`
}

type ApiError struct {
	errorDto ErrorDto
}

func (ae ApiError) Error() string {
	return fmt.Sprintf("Received an error from API. code: %d. message: %s", ae.errorDto.Code, ae.errorDto.ErrorInfo)
}

func DecodeError(r *http.Response) error {
	errorDto, err := Decode[ErrorDto](r)
	if err != nil {
		return err
	}
	return ApiError{errorDto: errorDto}
}

func Decode[T any](r *http.Response) (T, error) {
	var v T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return v, err
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return v, fmt.Errorf("failed to decode json: %w", err)
	}
	return v, nil
}
