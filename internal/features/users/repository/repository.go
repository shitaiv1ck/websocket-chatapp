package users_repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

type UsersRepository struct {
	store core_postgres.Store
}

func NewRepository(store core_postgres.Store) *UsersRepository {
	return &UsersRepository{
		store: store,
	}
}

func (r *UsersRepository) Save(ctx context.Context, user domains.User) (domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO chat.users(username, password_hash)
		VALUES($1, $2)
		RETURNING id, username;
	`

	var createdUser domains.User
	if err := r.store.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
	).Scan(
		&createdUser.ID,
		&createdUser.Username,
	); err != nil {
		if errPQ, ok := err.(*pgconn.PgError); ok {
			if errPQ.Code == "23505" {
				return domains.User{}, core_errors.ErrConflict
			}
		}

		return domains.User{}, err
	}

	return createdUser, nil
}

func (r *UsersRepository) FindAll(ctx context.Context, limit *int, offset *int) ([]domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT id, username, is_online
		FROM chat.users
		LIMIT $1
		OFFSET $2;
	`

	rows, err := r.store.Query(ctx, query, limit, offset)
	if err != nil {
		return []domains.User{}, err
	}
	defer rows.Close()

	users := make([]domains.User, 0)
	for rows.Next() {
		var user domains.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.IsOnline,
		); err != nil {
			return []domains.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UsersRepository) FindByID(ctx context.Context, id int) (domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM chat.users
		WHERE id = $1;
	`

	var user domains.User
	if err := r.store.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.IsOnline,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.User{}, core_errors.ErrNotFound
		}

		return domains.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) FindByUsername(ctx context.Context, username string) (domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM chat.users
		WHERE username = $1;
	`

	var user domains.User
	if err := r.store.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.IsOnline,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.User{}, core_errors.ErrNotFound
		}

		return domains.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) Update(ctx context.Context, user domains.User) (domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		UPDATE chat.users
		SET username = $1, password_hash = $2
		WHERE id = $3
		RETURNING id, username, is_online;
	`

	var patchedUser domains.User
	if err := r.store.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
		user.ID,
	).Scan(
		&patchedUser.ID,
		&patchedUser.Username,
		&patchedUser.IsOnline,
	); err != nil {
		if errPQ, ok := err.(*pgconn.PgError); ok {
			if errPQ.Code == "23505" {
				return domains.User{}, core_errors.ErrConflict
			}
		}

		return domains.User{}, err
	}

	return patchedUser, nil
}
