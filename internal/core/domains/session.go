package domains

import (
	"fmt"
	"time"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type Session struct {
	SessionToken string
	CSRFToken    string
	UserID       int
	CreatedAt    time.Time
	ExpiredAt    time.Time
}

func (s *Session) Validate() error {
	if s.UserID <= 0 {
		return fmt.Errorf("user id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if !s.ExpiredAt.After(time.Now()) {
		return fmt.Errorf("expired token: %w", core_errors.ErrCoockie)
	}

	return nil
}
