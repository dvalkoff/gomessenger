package user

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type ContactsService interface {
	AddContact(userId uuid.UUID, contactId uuid.UUID) error
	GetContacts(userId uuid.UUID) ([]UserInfo, error)
}

type contactsService struct {
	userRepository UserRepository
}

func NewContactsService(userRepository UserRepository) ContactsService {
	return &contactsService{userRepository: userRepository}
}

func (service *contactsService) AddContact(userId uuid.UUID, contactId uuid.UUID) error {
	if userId == contactId {
		slog.Error("The user cannot add himself to the contacts list", "userId", userId)
		return fmt.Errorf("The user cannot add himself to the contacts list")
	}
	return service.userRepository.AddContact(userId, contactId)
}

func (service *contactsService) GetContacts(userId uuid.UUID) ([]UserInfo, error) {
	userRows, err := service.userRepository.GetContacts(userId)
	if err != nil {
		return nil, err
	}
	users := make([]UserInfo, len(userRows))
	for i, userRow := range userRows {
		users[i] = UserInfo{
			Id:       userRow.Id,
			Nickname: userRow.Nickname,
		}
	}
	return users, nil
}
