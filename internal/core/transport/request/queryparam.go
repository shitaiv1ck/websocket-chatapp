package core_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert '%v' to number: %w", value, core_errors.ErrInvalidArg)
	}

	return &num, nil
}

func GetStringQueryParam(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}

	return &value
}
