package messages_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

type MessagesRepository struct {
	store core_postgres.Store
}

func NewRepository(store core_postgres.Store) *MessagesRepository {
	return &MessagesRepository{
		store: store,
	}
}

func (r *MessagesRepository) Save(ctx context.Context, message domains.Message) (domains.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		WITH inserted AS (
			INSERT INTO chat.messages(chat_id, sender_id, receiver_id, content)
			VALUES($1, $2, $3, $4)
			RETURNING *
		)
		UPDATE chat.chats AS c
		SET last_message_content = i.content, last_message_at = i.created_at
		FROM inserted AS i
		WHERE c.id = i.chat_id
		AND c.user1_id = LEAST(i.sender_id, i.receiver_id)
		AND c.user2_id = GREATEST(i.sender_id, i.receiver_id)
		RETURNING i.id, i.chat_id, i.sender_id, i.receiver_id, i.content, i.created_at;
	`

	var createdMessage domains.Message
	if err := r.store.QueryRow(
		ctx,
		query,
		message.ChatID,
		message.SenderID,
		message.ReceiverID,
		message.Content,
	).Scan(
		&createdMessage.ID,
		&createdMessage.ChatID,
		&createdMessage.SenderID,
		&createdMessage.ReceiverID,
		&createdMessage.Content,
		&createdMessage.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Message{}, fmt.Errorf(
				"chat with id=%v isn'n your chat with user with id=%v: %w",
				message.ChatID,
				message.ReceiverID,
				core_errors.ErrNotFound,
			)
		}

		if pgxError, ok := err.(*pgconn.PgError); ok {
			if pgxError.Code == "23503" {
				return domains.Message{}, fmt.Errorf(
					"chat with id=%v doesn't exist: %w",
					message.ChatID,
					core_errors.ErrNotFound,
				)
			}
		}

		return domains.Message{}, err
	}

	return createdMessage, nil
}

func (r *MessagesRepository) FindByChatIDAndUserID(ctx context.Context, userID, chatID int) ([]domains.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM chat.messages
		WHERE chat_id = $1 AND (sender_id = $2 OR receiver_id = $2)
		ORDER BY created_at ASC;
	`

	rows, err := r.store.Query(ctx, query, chatID, userID)
	if err != nil {
		return []domains.Message{}, err
	}
	defer rows.Close()

	messages := make([]domains.Message, 0)
	for rows.Next() {
		var message domains.Message
		if err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Content,
			&message.CreatedAt,
		); err != nil {
			return []domains.Message{}, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}
