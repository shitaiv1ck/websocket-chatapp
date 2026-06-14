package domains

import "time"

type Chat struct {
	ID                 int
	FirstUser          User
	SecondUser         User
	LastMessageContent *string
	LastMessageAt      time.Time
}
