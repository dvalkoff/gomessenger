package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

const (
	jsonContentType = "application/json"
)

type UserRegistrationInfo struct {
	Nickname          string   `json:"nickname"`
	Password          string   `json:"password"`
	IdentityPubKey    string   `json:"identityPubKey"`
	SignedPubPreKey   string   `json:"signedPubPrekey"`
	OneTimePubPreKeys []string `json:"oneTimePubPrekeys"`
}

type UserRegistratedInfo struct {
	Id uuid.UUID `json:"id"`
}

type UserAuthenticationInfo struct {
	UserId   uuid.UUID `json:"userId"`
	Password string    `json:"password"`
}

type AuthenticationInto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserApi interface {
	RegisterUser(baseURL string, dto UserRegistrationInfo) (UserRegistratedInfo, error)
	LogIn(baseURL string, dto UserAuthenticationInfo) (AuthenticationInto, error)
}

type userApi struct {
}

func NewUserApi() UserApi {
	return &userApi{}
}

func (api *userApi) RegisterUser(baseURL string, dto UserRegistrationInfo) (UserRegistratedInfo, error) {
	url, err := url.JoinPath(baseURL, "/signup")
	if err != nil {
		return UserRegistratedInfo{}, err
	}
	requestBody, err := json.Marshal(dto)
	if err != nil {
		return UserRegistratedInfo{}, err
	}
	response, err := http.Post(url, jsonContentType, bytes.NewReader(requestBody))
	if err != nil {
		return UserRegistratedInfo{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return UserRegistratedInfo{}, DecodeError(response)
	}
	userRegistratedDto, err := Decode[UserRegistratedInfo](response)
	if err != nil {
		return UserRegistratedInfo{}, err
	}
	return userRegistratedDto, nil
}

func (api *userApi) LogIn(baseURL string, dto UserAuthenticationInfo) (AuthenticationInto, error) {
	return AuthenticationInto{}, nil
}
