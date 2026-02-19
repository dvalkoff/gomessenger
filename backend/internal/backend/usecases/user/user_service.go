package user

import (
	"fmt"
	"log/slog"

	"encoding/base64"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	OneTimePreKeysAmount = 100

	FindUsersByNicknamePrefixMinLen = 4
)

type UserService interface {
	RegisterUser(UserRegistrationInfo) (UserRegistratedInfo, error)
	FindUsers(nicknameSubstring string) ([]UserInfo, error)
}

type userService struct {
	userRepository   UserRepository
	pubKeyRepository PubKeyRepository
}

func NewUserService(userRepository UserRepository, pubKeyRepository PubKeyRepository) UserService {
	return &userService{userRepository: userRepository, pubKeyRepository: pubKeyRepository}
}

func (service *userService) RegisterUser(userDto UserRegistrationInfo) (UserRegistratedInfo, error) {
	userId, err := uuid.NewRandom()
	if err != nil {
		slog.Error("Failed to generate uuid", "error", err)
	}

	identityKey, err := base64.StdEncoding.DecodeString(userDto.IdentityPubKey)
	if err != nil {
		slog.Error("Failed to decode identity key", "error", err)
		return UserRegistratedInfo{}, err
	}
	signedPubKey, err := base64.StdEncoding.DecodeString(userDto.SignedPubPreKey)
	if err != nil {
		slog.Error("Failed to decode signed pre key", "error", err)
		return UserRegistratedInfo{}, err
	}
	if len(userDto.OneTimePubPreKeys) != OneTimePreKeysAmount {
		slog.Error("One time pre key amount does not meet requirements")
		return UserRegistratedInfo{}, fmt.Errorf("One time pre key amount does not meet requirements")
	}
	oneTimePreKeyRows := make([]PubKeyRow, 0, len(userDto.OneTimePubPreKeys))
	for _, encodedKey := range userDto.OneTimePubPreKeys {
		key, err := base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			slog.Error("Failed to decode signed pre key", "error", err)
			return UserRegistratedInfo{}, err
		}
		oneTimePreKeyRows = append(oneTimePreKeyRows, PubKeyRow{
			Key:    key,
			Type:   OneTimePrekey,
			UserId: userId,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return UserRegistratedInfo{}, err
	}

	row := UserRow{
		Id:              userId,
		Nickname:        userDto.Nickname,
		HashedPassword:  hashedPassword,
		IdentityPubKey:  identityKey,
		SignedPubPrekey: signedPubKey,
	}

	tx, err := service.userRepository.StartTx()
	defer tx.Rollback()
	if err != nil {
		slog.Error("Failed to start transaction", "error", err)
		return UserRegistratedInfo{}, err
	}
	err = service.userRepository.SaveUser(tx, row)
	if err != nil {
		slog.Error("Failed to save user", "error", err)
		return UserRegistratedInfo{}, err
	}
	err = service.pubKeyRepository.SaveKeys(tx, oneTimePreKeyRows)
	if err != nil {
		slog.Error("Failed to save keys", "error", err)
		return UserRegistratedInfo{}, err
	}
	return UserRegistratedInfo{Id: userId}, tx.Commit()
}

func (service *userService) FindUsers(nicknameSubstring string) ([]UserInfo, error) {
	if len(nicknameSubstring) < FindUsersByNicknamePrefixMinLen {
		slog.Error("Nickname prefix is too short")
		return nil, fmt.Errorf("Nickname prefix is too short")
	}
	users, err := service.userRepository.FindUsersByNickname(nicknameSubstring)
	if err != nil {
		return nil, err
	}
	mappedUsers := make([]UserInfo, 0, len(users))
	for _, user := range users {
		mappedUsers = append(mappedUsers, UserInfo{
			Id:       user.Id,
			Nickname: user.Nickname,
		})
	}
	return mappedUsers, nil
}
