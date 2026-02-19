package utils

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	UserContextKey = "user"
	UserIdKey      = "sub"
	ExpirationTime = "exp"
)

func GetUserId(ctx context.Context) uuid.UUID {
	user := ctx.Value(UserContextKey).(jwt.MapClaims)
	return user[UserIdKey].(uuid.UUID)
}
