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

func (rw *ResponseWriter) JsonResponse(body any, statusCode int) {
	rw.WriteHeader(statusCode)

	if err := json.NewEncoder(rw).Encode(body); err != nil {
		panic(err)
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

func (rw *ResponseWriter) GetStatusCode() int {
	return rw.statusCode
}

func errStatusCode(err error) int {
	if errors.Is(err, core_errors.ErrCoockie) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, core_errors.ErrInvalidArg) {
		return http.StatusBadRequest
	}

	if errors.Is(err, core_errors.ErrConflict) {
		return http.StatusConflict
	}

	if errors.Is(err, core_errors.ErrNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(err, core_errors.ErrUnauthenticate) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, core_errors.ErrUnauthorize) {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
