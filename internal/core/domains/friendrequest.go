package domains

import "time"

type FriendRequest struct {
	ID        int
	FromUser  UserBrief
	ToUser    UserBrief
	CreatedAt time.Time
}

type UserBrief struct {
	ID       int
	Username string
}
