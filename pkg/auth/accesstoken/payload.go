package accesstoken

import (
	"time"
)

type JwtPayload struct {
	UserId    int
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(userId int, ttl time.Duration) *JwtPayload {
	return &JwtPayload{
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(ttl),
	}
}
