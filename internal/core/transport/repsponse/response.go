package core_repsponse

import (
	"encoding/json"
	"errors"
	"net/http"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type ResponseWriter struct {
	http.ResponseWriter

	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

func (rw *ResponseWriter) ErrorResponse(err error, msg string) {
	statusCode := errStatusCode(err)
	rw.WriteHeader(statusCode)

	errorDTO := ErrorDTO{
		Error:   err.Error(),
		Message: msg,
	}

	if err := json.NewEncoder(rw).Encode(errorDTO); err != nil {
		panic(err)
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)

	rw.statusCode = statusCode
}

func errStatusCode(err error) int {
	if errors.Is(err, core_errors.ErrCoockie) {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
