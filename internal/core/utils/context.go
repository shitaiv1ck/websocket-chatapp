package core_utils

import (
	"context"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

func GetIntFromContext(ctx context.Context, key string) (*int, error) {
	value := ctx.Value(key)

	if value == nil {
		return nil, core_errors.ErrInvalidArg
	}

	num, ok := value.(int)
	if !ok {
		return nil, core_errors.ErrInvalidArg
	}

	return &num, nil
}
