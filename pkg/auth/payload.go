package auth

import (
	"github.com/vaberof/ssugt-projects/pkg/domain"
	"time"
)

type JwtPayload struct {
	UserId    domain.UserId
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(userId domain.UserId, ttl time.Duration) *JwtPayload {
	return &JwtPayload{
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(ttl),
	}
}
