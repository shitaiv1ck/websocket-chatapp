package chats_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

type ChatsRepository struct {
	store core_postgres.Store
}

func NewRepository(store core_postgres.Store) *ChatsRepository {
	return &ChatsRepository{
		store: store,
	}
}

func (r *ChatsRepository) SaveOrFind(ctx context.Context, firstUserID int, secondUserID int) (domains.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		WITH inserted AS (
			INSERT INTO chat.chats(user1_id, user2_id)
			VALUES (LEAST($1::int, $2::int), GREATEST($1::int, $2::int))
			ON CONFLICT (user1_id, user2_id)
			DO UPDATE SET last_message_at = chat.chats.last_message_at
			RETURNING *
		)
		SELECT i.id,
			i.user1_id, u1.username, u1.is_online,
			i.user2_id, u2.username, u2.is_online,
			i.last_message_content,
			i.last_message_at
		FROM inserted AS i
		JOIN chat.users AS u1 ON i.user1_id = u1.id
		JOIN chat.users AS u2 ON i.user2_id = u2.id;
	`

	var chat domains.Chat
	if err := r.store.QueryRow(
		ctx,
		query,
		firstUserID,
		secondUserID,
	).Scan(
		&chat.ID,
		&chat.FirstUser.ID,
		&chat.FirstUser.Username,
		&chat.FirstUser.IsOnline,
		&chat.SecondUser.ID,
		&chat.SecondUser.Username,
		&chat.SecondUser.IsOnline,
		&chat.LastMessageContent,
		&chat.LastMessageAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Chat{}, core_errors.ErrNotFound
		}

		return domains.Chat{}, err
	}

	return chat, nil
}

func (r *ChatsRepository) FindByUserID(ctx context.Context, userID int, limit *int, offset *int) ([]domains.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT c.id,
			c.user1_id, u1.username, u1.is_online,
			c.user2_id, u2.username, u2.is_online,
			c.last_message_content, c.last_message_at
		FROM chat.chats AS c
		JOIN chat.users AS u1 ON c.user1_id = u1.id
		JOIN chat.users AS u2 ON c.user2_id = u2.id
		WHERE c.user1_id = $1 OR c.user2_id = $1
		ORDER BY c.last_message_at DESC
		LIMIT $2
		OFFSET $3;
	`

	rows, err := r.store.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return []domains.Chat{}, err
	}
	defer rows.Close()

	chats := make([]domains.Chat, 0)
	for rows.Next() {
		var chat domains.Chat
		if err := rows.Scan(
			&chat.ID,
			&chat.FirstUser.ID,
			&chat.FirstUser.Username,
			&chat.FirstUser.IsOnline,
			&chat.SecondUser.ID,
			&chat.SecondUser.Username,
			&chat.SecondUser.IsOnline,
			&chat.LastMessageContent,
			&chat.LastMessageAt,
		); err != nil {
			return []domains.Chat{}, err
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (r *ChatsRepository) Delete(ctx context.Context, userID int, chatID int) (domains.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		DELETE FROM chat.chats
		WHERE id = $1 AND (user1_id = $2 OR user2_id = $2)
		RETURNING id, user1_id, user2_id;
	`

	var deletedChat domains.Chat
	if err := r.store.QueryRow(
		ctx,
		query,
		chatID,
		userID,
	).Scan(
		&deletedChat.ID,
		&deletedChat.FirstUser.ID,
		&deletedChat.SecondUser.ID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Chat{}, fmt.Errorf("chat doesn't exist: %w", core_errors.ErrNotFound)
		}

		return domains.Chat{}, err
	}

	return deletedChat, nil
}
