package model

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/dvalkoff/gomessenger/gui/internal/events"
	"github.com/dvalkoff/gomessenger/gui/internal/integration/api"
	"github.com/dvalkoff/gomessenger/gui/internal/integration/repository"
	"github.com/dvalkoff/gomessenger/gui/internal/utils"
	"github.com/google/uuid"
)

const (
	oneTimePrekeyAmount = 100
)

type User struct {
	Id       uuid.UUID
	Nickname string

	AccessToken  string
	RefreshToken string

	IdentityKey    *Key
	SignedPrekeys  []*Key
	OneTimePrekeys []*Key
}

type Key = ecdh.PrivateKey

var currentUser *User
var userLock *sync.RWMutex = &sync.RWMutex{}

func CurrentUser() *User {
	userLock.RLock()
	defer userLock.RUnlock()

	return currentUser
}

func setCurrentUser(newUser *User) {
	userLock.Lock()
	defer userLock.Unlock()

	currentUser = newUser
}

type UserService interface {
	InitUser(events.InitUserEvent)
	CreateUser(events.UserSignUpAttemptedEvent)
}

type userService struct {
	userApi        api.UserApi
	userRepository repository.UserRepository
	eventStream    chan<- events.Event
}

func NewUserService(
	userApi api.UserApi,
	userRepository repository.UserRepository,
	eventStream chan<- events.Event,
) UserService {
	return &userService{
		userApi:        userApi,
		userRepository: userRepository,
		eventStream:    eventStream,
	}
}

func (s *userService) InitUser(event events.InitUserEvent) {
	workspace := CurrentWorkspace()
	userRow, err := s.userRepository.GetCurrentUser(workspace.Id)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	if userRow == nil {
		s.eventStream <- events.SwitchViewEvent{NewView: events.SignUpView}
		return
	}
	if len(userRow.AccessToken) == 0 && len(userRow.RefreshToken) == 0 {
		s.eventStream <- events.SwitchViewEvent{NewView: events.SignInView}
		return
	}
	identityKeys := userRow.GetKeys(repository.IdentityKey)
	if len(identityKeys) != 1 {
		s.eventStream <- events.ErrorNotificationEvent{Err: fmt.Errorf("Identity keys count is invalid")}
		return
	}
	mappedIdentityKey, err := utils.MapKey(identityKeys[0].PrivateKey)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	mappedSignedPrekeys, err := mapKeys(userRow.GetKeys(repository.SignedPrekey))
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	mappedOneTimePrekeys, err := mapKeys(userRow.GetKeys(repository.OneTimePrekey))
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	user := &User{
		Id:             userRow.Id,
		Nickname:       userRow.Nickname,
		AccessToken:    userRow.AccessToken,
		RefreshToken:   userRow.RefreshToken,
		IdentityKey:    mappedIdentityKey,
		SignedPrekeys:  mappedSignedPrekeys,
		OneTimePrekeys: mappedOneTimePrekeys,
	}

	setCurrentUser(user)

	s.eventStream <- events.SwitchViewEvent{NewView: events.MainAppView}
}

func mapKeys(rows []repository.KeyRow) ([]*Key, error) {
	result := make([]*Key, 0, len(rows))
	for _, row := range rows {
		key, err := utils.MapKey(row.PrivateKey)
		if err != nil {
			return nil, err
		}
		result = append(result, key)
	}
	return result, nil
}

func (s *userService) CreateUser(event events.UserSignUpAttemptedEvent) {
	identityKey, err := utils.GenerateKey()
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	signedPrekeys, err := utils.GenerateKeys(1)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	oneTimePrekeys, err := utils.GenerateKeys(oneTimePrekeyAmount)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}

	userRegistratedInfo, err := s.userApi.RegisterUser(
		CurrentWorkspace().URL,
		api.UserRegistrationInfo{
			Nickname:          event.Nickname,
			Password:          event.Password,
			IdentityPubKey:    EncodePublicKey(identityKey),
			SignedPubPreKey:   EncodePublicKey(signedPrekeys[0]),
			OneTimePubPreKeys: EncodePublicKeys(oneTimePrekeys),
		},
	)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}

	user := &User{}
	user.Id = userRegistratedInfo.Id
	user.Nickname = event.Nickname
	user.IdentityKey = identityKey
	user.SignedPrekeys = signedPrekeys
	user.OneTimePrekeys = oneTimePrekeys

	err = s.SaveUser(user)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	setCurrentUser(user)
	s.eventStream <- events.SwitchViewEvent{NewView: events.SignInView}
}

func EncodePublicKeys(keys []*ecdh.PrivateKey) []string {
	encodedSlice := make([]string, 0, len(keys))
	for _, key := range keys {
		encodedSlice = append(encodedSlice, EncodePublicKey(key))
	}
	return encodedSlice
}

func EncodePublicKey(key *ecdh.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(key.PublicKey().Bytes())
}

func (s *userService) SaveUser(user *User) error {
	keyRows := make([]repository.KeyRow, 0, len(user.SignedPrekeys)+len(user.OneTimePrekeys)+1)
	keyRows = append(keyRows,
		repository.KeyRow{
			UserId:     user.Id,
			Type:       repository.IdentityKey,
			PrivateKey: user.IdentityKey.Bytes(),
		},
	)
	for _, key := range user.SignedPrekeys {
		keyRows = append(keyRows,
			repository.KeyRow{
				UserId:     user.Id,
				Type:       repository.SignedPrekey,
				PrivateKey: key.Bytes(),
			},
		)
	}
	for _, key := range user.OneTimePrekeys {
		keyRows = append(keyRows,
			repository.KeyRow{
				UserId:     user.Id,
				Type:       repository.OneTimePrekey,
				PrivateKey: key.Bytes(),
			},
		)
	}
	userRow := repository.UserRow{
		Id:          user.Id,
		WorkspaceId: CurrentWorkspace().Id,
		Nickname:    user.Nickname,
		Keys:        keyRows,
	}
	err := s.userRepository.SaveUser(userRow)
	if err != nil {
		return err
	}
	return s.userRepository.SetUserCurrent(user.Id)
}
