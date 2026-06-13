package domains

import (
	"fmt"
	"time"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type FriendRequest struct {
	ID        int
	FromUser  User
	ToUser    User
	CreatedAt time.Time
}

func (fr *FriendRequest) Validate() error {
	if fr.FromUser.ID <= 0 || fr.ToUser.ID <= 0 {
		return fmt.Errorf("user id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if fr.FromUser.ID == fr.ToUser.ID {
		return fmt.Errorf("FromUserId can't be equal ToUserID: %w", core_errors.ErrInvalidArg)
	}

	return nil
}
