package repository

import "github.com/google/uuid"

type ChatRow struct {
	Id uuid.UUID
}

type ChatRepository interface {
	GetChats(userId uuid.UUID) ([]ChatRow, error)
}

func NewChatRepository() ChatRepository {
	return nil
}
