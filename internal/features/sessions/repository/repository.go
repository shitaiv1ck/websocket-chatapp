package sessions_repository

import (
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

type SessionsRepository struct {
	store *core_postgres.Store
}

func NewRepository(store *core_postgres.Store) *SessionsRepository {
	return &SessionsRepository{
		store: store,
	}
}

func (r *SessionsRepository) Save(session domains.Session) (domains.Session, error) {
	query := `
		INSERT INTO chat.sessions(session_token, csrf_token, user_id, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *;
	`

	var createdSession domains.Session
	if err := r.store.QueryRow(
		query,
		&session.SessionToken,
		&session.CSRFToken,
		&session.UserID,
		&session.ExpiresAt,
	).Scan(
		&createdSession.SessionToken,
		&createdSession.CSRFToken,
		&createdSession.UserID,
		&createdSession.CreatedAt,
		&createdSession.ExpiresAt,
	); err != nil {
		return domains.Session{}, err
	}

	return createdSession, nil
}

func (r *SessionsRepository) FindByToken(token string) (domains.Session, error) {
	query := `
		SELECT * FROM chat.sessions
		WHERE session_token = $1;
	`

	var session domains.Session
	if err := r.store.QueryRow(
		query,
		token,
	).Scan(
		&session.SessionToken,
		&session.CSRFToken,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	); err != nil {
		return domains.Session{}, err
	}

	return session, nil
}

func (r *SessionsRepository) Delete(token string) error {
	query := `
		DELETE FROM chat.sessions
		WHERE session_token = $1;
	`

	if _, err := r.store.Exec(
		query,
		token,
	); err != nil {
		return err
	}

	return nil
}
