package helper

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserContextKey = "user"
	UserNicknameKey = "nickname"
)

func GetNickname(ctx context.Context) string {
	user := ctx.Value(UserContextKey).(jwt.MapClaims)
	return user[UserNicknameKey].(string)
}
