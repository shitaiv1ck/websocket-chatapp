package domains

import (
	"fmt"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type Friendship struct {
	ID         int
	FirstUser  User
	SecondUser User
}

func (f *Friendship) Validate() error {
	if f.FirstUser.ID <= 0 || f.SecondUser.ID <= 0 {
		return fmt.Errorf("user id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if f.FirstUser.ID == f.SecondUser.ID {
		return fmt.Errorf("first user id can't be equal second user id: %w", core_errors.ErrInvalidArg)
	}

	return nil
}
