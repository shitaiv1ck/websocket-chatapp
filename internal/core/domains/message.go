package domains

import (
	"fmt"
	"time"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type Message struct {
	ID         int
	ChatID     int
	SenderID   int
	ReceiverID int
	Content    string
	CreatedAt  time.Time
}

func (m *Message) Validate() error {
	if m.ChatID <= 0 {
		return fmt.Errorf("chat id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if m.SenderID <= 0 {
		return fmt.Errorf("sender id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if m.ReceiverID <= 0 {
		return fmt.Errorf("receiver id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if m.SenderID == m.ReceiverID {
		return fmt.Errorf("can't send message to yourself: %w", core_errors.ErrInvalidArg)
	}

	if len([]rune(m.Content)) == 0 {
		return fmt.Errorf("message can't be empty")
	}

	return nil
}
