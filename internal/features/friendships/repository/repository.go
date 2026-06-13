package friendships_repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

type FriendshipRepository struct {
	store *core_postgres.Store
}

func NewRepository(store *core_postgres.Store) *FriendshipRepository {
	return &FriendshipRepository{
		store: store,
	}
}

func (r *FriendshipRepository) Save(fromUserID int, toUserID int) (domains.Friendship, error) {
	query := `
		WITH inserted AS (
			INSERT INTO chat.friendships(user1_id, user2_id)
			VALUES (LEAST($1::int, $2::int), GREATEST($1::int, $2::int))
			RETURNING *
		)
		SELECT i.id,
			i.user1_id, u1.username, u1.is_online,
			i.user2_id, u2.username, u2.is_online
		FROM inserted AS i
		JOIN chat.users AS u1 ON i.user1_id = u1.id
		JOIN chat.users AS u2 ON i.user2_id = u2.id;
	`

	var createdFriendship domains.Friendship
	if err := r.store.QueryRow(
		query,
		fromUserID,
		toUserID,
	).Scan(
		&createdFriendship.ID,
		&createdFriendship.FirstUser.ID,
		&createdFriendship.FirstUser.Username,
		&createdFriendship.FirstUser.IsOnline,
		&createdFriendship.SecondUser.ID,
		&createdFriendship.SecondUser.Username,
		&createdFriendship.SecondUser.IsOnline,
	); err != nil {
		if errPQ, ok := err.(*pq.Error); ok {
			if errPQ.Code == "23505" {
				return domains.Friendship{}, core_errors.ErrNotFound
			}
		}

		return domains.Friendship{}, err
	}

	return createdFriendship, nil
}

func (r *FriendshipRepository) FindByUsers(firstUserID int, secondUserID int) (domains.Friendship, error) {
	query := `
		SELECT * FROM chat.friendships
		WHERE user1_id = LEAST($1::int, $2::int) AND user2_id = GREATEST($1::int, $2::int);
	`

	var foundFriendship domains.Friendship
	if err := r.store.QueryRow(
		query,
		firstUserID,
		secondUserID,
	).Scan(
		&foundFriendship.ID,
		&foundFriendship.FirstUser.ID,
		&foundFriendship.SecondUser.ID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domains.Friendship{}, core_errors.ErrNotFound
		}

		return domains.Friendship{}, err
	}

	return foundFriendship, nil
}

func (r *FriendshipRepository) FindByUserID(userID int, limit *int, offset *int) ([]domains.Friendship, error) {
	query := `
		SELECT f.id,
			f.user1_id, u1.username, u1.is_online,
			f.user2_id, u2.username, u2.is_online
		FROM chat.friendships AS f
		JOIN chat.users AS u1 ON f.user1_id = u1.id
		JOIN chat.users AS u2 ON f.user2_id = u2.id
		WHERE f.user1_id = $1 OR f.user2_id = $1
		LIMIT $2
		OFFSET $3;
	`

	rows, err := r.store.Query(query, userID, limit, offset)
	if err != nil {
		return []domains.Friendship{}, err
	}
	defer rows.Close()

	friendships := make([]domains.Friendship, 0)
	for rows.Next() {
		var friendship domains.Friendship
		if err := rows.Scan(
			&friendship.ID,
			&friendship.FirstUser.ID,
			&friendship.FirstUser.Username,
			&friendship.FirstUser.IsOnline,
			&friendship.SecondUser.ID,
			&friendship.SecondUser.Username,
			&friendship.SecondUser.IsOnline,
		); err != nil {
			return []domains.Friendship{}, err
		}

		friendships = append(friendships, friendship)
	}

	return friendships, nil
}

func (r *FriendshipRepository) Delete(userID int, friendshipID int) error {
	query := `
		DELETE FROM chat.friendships
		WHERE id = $1 AND (user1_id = $2 OR user2_id = $2);
	`

	result, err := r.store.Exec(query, friendshipID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("friendship doesn't exist: %w", core_errors.ErrNotFound)
	}

	return nil
}
