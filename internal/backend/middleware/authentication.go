package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dvalkoff/gomessenger/internal/backend/helper"
	"github.com/dvalkoff/gomessenger/internal/backend/usecases/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticationInfo struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type AuthenticationInto struct {
	AccessToken string `json:"accessToken"`
}

type AuthenticationProvider interface {
	LogIn() http.Handler
	AuthMiddleware(next http.Handler) http.Handler
}

type authenticationProvider struct {
	userRepository user.UserRepository
	signingJwtSecret []byte
}

func NewAuthenticationProvider(userRepository user.UserRepository, signingJwtSecret string) AuthenticationProvider {
	return &authenticationProvider{
		userRepository: userRepository, 
		signingJwtSecret: []byte(signingJwtSecret),
	}
}

func (authProvider *authenticationProvider) LogIn() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authInfo, err := helper.Decode[UserAuthenticationInfo](r)

			users, err := authProvider.userRepository.FindUsersByNickname(authInfo.Nickname)
			if err != nil || len(users) == 0 {
				slog.Error("Failed to get user by nickname", "nickname", authInfo.Nickname, "error", err)
				err := fmt.Errorf("User not found %w", err)
				helper.EncodeError(w, r, http.StatusUnauthorized, err)
				return
			}
			user := users[0]
			err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(authInfo.Password))
			if err != nil {
				slog.Error("Password is not correct", "nickname", authInfo.Nickname, "error", err)
				helper.EncodeError(w, r, http.StatusUnauthorized, err)
				return
			}

			jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				helper.UserNicknameKey: authInfo.Nickname,
			})
			tokenStr, err := jwt.SignedString(authProvider.signingJwtSecret)
			if err != nil {
				slog.Error("Failed to sign token", "error", err)
				helper.EncodeError(w, r, http.StatusUnauthorized, err)
				return
			}
			
			authDetails := AuthenticationInto{AccessToken: tokenStr}
			err = helper.Encode(w, r, http.StatusOK, authDetails)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}


func (authProvider *authenticationProvider) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing auth header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return authProvider.signingJwtSecret, nil
			})

			if err != nil || !token.Valid {
				helper.EncodeError(w, r, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			ctx := context.WithValue(r.Context(), helper.UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
