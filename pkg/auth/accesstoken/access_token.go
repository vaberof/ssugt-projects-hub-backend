package accesstoken

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"ssugt-projects-hub/config"
	"strconv"
	"time"
)

var (
	ErrInvalidToken         = errors.New("token is invalid")
	ErrInvalidSigningMethod = errors.New("signing method is invalid")
	ErrExpiredToken         = errors.New("token has expired")
)

func Create(userId int, ttl time.Duration) (string, error) {
	payload := NewPayload(userId, ttl)

	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(userId),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
	})

	token, err := jwtWithClaims.SignedString([]byte(config.SecretKey()))

	return token, err
}

func Verify(token string) (*JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(config.SecretKey()), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if hasExpired(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}

	userId, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload := &JwtPayload{
		UserId:    userId,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}

	return payload, nil
}

func hasExpired(expireTime time.Time) bool {
	currentTime := time.Now().UTC()
	return currentTime.After(expireTime.UTC())
}
