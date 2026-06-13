package friendrequests_respository

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

type FriendRequestsRepository struct {
	store core_postgres.Store
}

func NewRepository(store core_postgres.Store) *FriendRequestsRepository {
	return &FriendRequestsRepository{
		store: store,
	}
}

func (r *FriendRequestsRepository) Save(ctx context.Context, request domains.FriendRequest) (domains.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		WITH inserted AS (
			INSERT INTO chat.friendrequests(from_id, to_id)
			VALUES($1, $2)
			RETURNING *
		)
		SELECT i.id,
			i.from_id, u1.username, u1.is_online,
			i.to_id, u2.username, u2.is_online,
			i.created_at
		FROM inserted AS i
		JOIN chat.users AS u1 ON i.from_id = u1.id
		JOIN chat.users AS u2 ON i.to_id = u2.id;
	`

	var createdFriendRequest domains.FriendRequest
	if err := r.store.QueryRow(
		ctx,
		query,
		request.FromUser.ID,
		request.ToUser.ID,
	).Scan(
		&createdFriendRequest.ID,
		&createdFriendRequest.FromUser.ID,
		&createdFriendRequest.FromUser.Username,
		&createdFriendRequest.FromUser.IsOnline,
		&createdFriendRequest.ToUser.ID,
		&createdFriendRequest.ToUser.Username,
		&createdFriendRequest.ToUser.IsOnline,
		&createdFriendRequest.CreatedAt,
	); err != nil {
		if errPQ, ok := err.(*pgconn.PgError); ok {
			if errPQ.Code == "23505" {
				return domains.FriendRequest{}, core_errors.ErrConflict
			}

			if errPQ.Code == "23503" {
				return domains.FriendRequest{}, fmt.Errorf("user with id=%v doesn't exist: %w", request.ToUser.ID, core_errors.ErrInvalidArg)
			}
		}

		return domains.FriendRequest{}, err
	}

	return createdFriendRequest, nil
}

func (r *FriendRequestsRepository) FindByFromIDAndToID(ctx context.Context, fromID int, toID int) (domains.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM chat.friendrequests
		WHERE from_id = $1 AND to_id = $2;
	`

	var foundFriendRequest domains.FriendRequest
	if err := r.store.QueryRow(
		ctx,
		query,
		fromID,
		toID,
	).Scan(
		&foundFriendRequest.ID,
		&foundFriendRequest.FromUser.ID,
		&foundFriendRequest.ToUser.ID,
		&foundFriendRequest.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.FriendRequest{}, core_errors.ErrNotFound
		}

		return domains.FriendRequest{}, err
	}

	return foundFriendRequest, nil
}

func (r *FriendRequestsRepository) FindByToID(ctx context.Context, toID int) ([]domains.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT r.id, r.from_id, ur.username, ur.is_online, r.to_id, ut.username, ut.is_online, r.created_at
		FROM chat.friendrequests AS r
		JOIN chat.users AS ur ON r.from_id = ur.id
		JOIN chat.users AS ut ON r.to_id = ut.id
		WHERE r.to_id = $1;
	`
	rows, err := r.store.Query(ctx, query, toID)
	if err != nil {
		return []domains.FriendRequest{}, err
	}
	defer rows.Close()

	friendRequests := make([]domains.FriendRequest, 0)
	for rows.Next() {
		var request domains.FriendRequest
		if err := rows.Scan(
			&request.ID,
			&request.FromUser.ID,
			&request.FromUser.Username,
			&request.FromUser.IsOnline,
			&request.ToUser.ID,
			&request.ToUser.Username,
			&request.ToUser.IsOnline,
			&request.CreatedAt,
		); err != nil {
			return []domains.FriendRequest{}, err
		}

		friendRequests = append(friendRequests, request)
	}

	return friendRequests, nil
}

func (r *FriendRequestsRepository) FindByIDAndToID(ctx context.Context, requestID int, toID int) (domains.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM chat.friendrequests
		WHERE id = $1 AND to_id = $2;
	`

	var request domains.FriendRequest
	if err := r.store.QueryRow(
		ctx,
		query,
		requestID,
		toID,
	).Scan(
		&request.ID,
		&request.FromUser.ID,
		&request.ToUser.ID,
		&request.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.FriendRequest{}, core_errors.ErrNotFound
		}

		return domains.FriendRequest{}, err
	}

	return request, nil
}

func (r *FriendRequestsRepository) FindByFromID(ctx context.Context, fromID int) ([]domains.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT r.id, r.from_id, ur.username, r.to_id, ut.username, r.created_at
		FROM chat.friendrequests AS r
		JOIN chat.users AS ur ON r.from_id = ur.id
		JOIN chat.users AS ut ON r.to_id = ut.id
		WHERE r.from_id = $1;
	`
	rows, err := r.store.Query(ctx, query, fromID)
	if err != nil {
		return []domains.FriendRequest{}, err
	}

	friendRequests := make([]domains.FriendRequest, 0)
	for rows.Next() {
		var request domains.FriendRequest
		if err := rows.Scan(
			&request.ID,
			&request.FromUser.ID,
			&request.FromUser.Username,
			&request.ToUser.ID,
			&request.ToUser.Username,
			&request.CreatedAt,
		); err != nil {
			return []domains.FriendRequest{}, err
		}

		friendRequests = append(friendRequests, request)
	}

	return friendRequests, nil
}

func (r *FriendRequestsRepository) Delete(ctx context.Context, userID int, requestID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		DELETE FROM chat.friendrequests
		WHERE id = $1 AND to_id = $2
	`

	result, err := r.store.Exec(ctx, query, requestID, userID)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows != 1 {
		return fmt.Errorf("friend request doesn't exist: %w", core_errors.ErrNotFound)
	}

	return nil
}

func (r *FriendRequestsRepository) DeleteTx(ctx context.Context, tx core_postgres.SQLExecuter, userID int, requestID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		DELETE FROM chat.friendrequests
		WHERE id = $1 AND to_id = $2
	`

	result, err := tx.Exec(ctx, query, requestID, userID)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows != 1 {
		return fmt.Errorf("friend request doesn't exist: %w", core_errors.ErrNotFound)
	}

	return nil
}
