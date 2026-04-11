package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, code int, message string, data any) {
	JSON(w, code, SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, message string, errors any) {
	JSON(w, code, ErrorResponse{
		Code:    code,
		Message: message,
		Errors:  errors,
	})
}
