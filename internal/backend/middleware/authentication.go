package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dvalkoff/gomessenger/internal/backend/usecases/user"
	"github.com/dvalkoff/gomessenger/internal/backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpirationTimeSeconds = 60 * 60 * 24 * 30
)

type UserAuthenticationInfo struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type AuthenticationInto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthenticationProvider interface {
	LogIn() http.Handler
	AuthMiddleware(next http.Handler) http.Handler
	AuthWsMiddleware(next http.Handler) http.Handler
}

type authenticationProvider struct {
	userRepository   user.UserRepository
	signingJwtSecret []byte
}

func NewAuthenticationProvider(userRepository user.UserRepository, signingJwtSecret string) AuthenticationProvider {
	return &authenticationProvider{
		userRepository:   userRepository,
		signingJwtSecret: []byte(signingJwtSecret),
	}
}

func (authProvider *authenticationProvider) LogIn() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authInfo, err := utils.Decode[UserAuthenticationInfo](r)

			users, err := authProvider.userRepository.FindUsersByNickname(authInfo.Nickname)
			if err != nil || len(users) == 0 {
				slog.Error("Failed to get user by nickname", "nickname", authInfo.Nickname, "error", err)
				err := fmt.Errorf("User not found %w", err)
				utils.EncodeError(w, r, http.StatusUnauthorized, err)
				return
			}
			user := users[0]
			err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(authInfo.Password))
			if err != nil {
				slog.Warn("Password is not correct", "nickname", authInfo.Nickname, "error", err)
				utils.EncodeError(w, r, http.StatusUnauthorized, err)
				return
			}

			tokenStr, err := authProvider.createAndSignToken(authInfo.Nickname)
			if err != nil {
				slog.Error("Failed to sign token", "error", err)
				utils.EncodeError(w, r, http.StatusUnauthorized, err)
				return
			}

			refreshToken := "" // TODO: implement refresh token

			authDetails := AuthenticationInto{
				AccessToken:  tokenStr,
				RefreshToken: refreshToken,
			}
			err = utils.Encode(w, r, http.StatusOK, authDetails)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}

func (authProvider *authenticationProvider) createAndSignToken(nickname string) (string, error) {
	expirationTime := time.Now().Add(accessTokenExpirationTimeSeconds * time.Second)
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		utils.UserNicknameKey: nickname,
		utils.ExpirationTime:  jwt.NewNumericDate(expirationTime),
	})
	return jwt.SignedString(authProvider.signingJwtSecret)
}

func (authProvider *authenticationProvider) AuthWsMiddleware(next http.Handler) http.Handler {
	return authProvider.authMiddleware(next, extractTokenFromQuery)
}

func (authProvider *authenticationProvider) AuthMiddleware(next http.Handler) http.Handler {
	return authProvider.authMiddleware(next, extractTokenFromHeaders)
}

func (authProvider *authenticationProvider) authMiddleware(next http.Handler, tokenExtractor func(r *http.Request) (string, error)) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := tokenExtractor(r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusUnauthorized, err)
			}
			token, err := jwt.Parse(
				tokenStr,
				func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method")
					}
					return authProvider.signingJwtSecret, nil
				},
			)
			if err != nil || !token.Valid || isTokenExpired(token) {
				utils.EncodeError(w, r, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			ctx := context.WithValue(r.Context(), utils.UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func extractTokenFromHeaders(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing auth header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("Missing auth header")
	}
	return parts[1], nil
}

func extractTokenFromQuery(r *http.Request) (string, error) {
	authToken := r.URL.Query().Get("token")
	if authToken == "" {
		return "", fmt.Errorf("Missing auth token")
	}
	return authToken, nil
}

func isTokenExpired(token *jwt.Token) bool {
	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return true
	}
	return expirationTime.Before(time.Now())
}
