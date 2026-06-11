package friendrequests_respository

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

type FriendRequestsRepository struct {
	store *core_postgres.Store
}

func NewRepository(store *core_postgres.Store) *FriendRequestsRepository {
	return &FriendRequestsRepository{
		store: store,
	}
}

func (r *FriendRequestsRepository) Save(request domains.FriendRequest) (domains.FriendRequest, error) {
	query := `
		INSERT INTO chat.friendrequests(from_id, to_id)
		VALUES($1, $2)
		RETURNING *;
	`

	var createdFriendRequest domains.FriendRequest
	if err := r.store.QueryRow(
		query,
		request.FromUser.ID,
		request.ToUser.ID,
	).Scan(
		&createdFriendRequest.ID,
		&createdFriendRequest.FromUser.ID,
		&createdFriendRequest.ToUser.ID,
		&createdFriendRequest.CreatedAt,
	); err != nil {
		if errPQ, ok := err.(*pq.Error); ok {
			if errPQ.Code == "23505" {
				return domains.FriendRequest{}, core_errors.ErrConflict
			}
		}

		return domains.FriendRequest{}, err
	}

	return createdFriendRequest, nil
}

func (r *FriendRequestsRepository) FindByFromIDAndToID(fromID int, toID int) (domains.FriendRequest, error) {
	query := `
		SELECT * FROM chat.friendrequests
		WHERE from_id = $1 AND to_id = $2;
	`

	var foundFriendRequest domains.FriendRequest
	if err := r.store.QueryRow(
		query,
		fromID,
		toID,
	).Scan(
		&foundFriendRequest.ID,
		&foundFriendRequest.FromUser.ID,
		&foundFriendRequest.ToUser.ID,
		&foundFriendRequest.CreatedAt,
	); err != nil {
		return domains.FriendRequest{}, err
	}

	return foundFriendRequest, nil
}

func (r *FriendRequestsRepository) FindByToID(toID int) ([]domains.FriendRequest, error) {
	query := `
		SELECT r.id, r.from_id, ur.username, r.to_id, ut.username, r.created_at
		FROM chat.friendrequests AS r
		JOIN chat.users AS ur ON r.from_id = ur.id
		JOIN chat.users AS ut ON r.to_id = ut.id
		WHERE r.to_id = $1;
	`
	rows, err := r.store.Query(query, toID)
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

func (r *FriendRequestsRepository) FindByFromID(fromID int) ([]domains.FriendRequest, error) {
	query := `
		SELECT r.id, r.from_id, ur.username, r.to_id, ut.username, r.created_at
		FROM chat.friendrequests AS r
		JOIN chat.users AS ur ON r.from_id = ur.id
		JOIN chat.users AS ut ON r.to_id = ut.id
		WHERE r.from_id = $1;
	`
	rows, err := r.store.Query(query, fromID)
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

func (r *FriendRequestsRepository) Delete(userID int, requestID int) error {
	query := `
		DELETE FROM chat.friendrequests
		WHERE id = $1 AND to_id = $2
	`

	result, err := r.store.Exec(query, requestID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("friend request doesn't exist: %w", core_errors.ErrNotFound)
	}

	return nil
}
