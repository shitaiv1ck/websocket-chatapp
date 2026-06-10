package domains

import "time"

type Session struct {
	SessionToken string
	CSRFToken    string
	UserID       int
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
