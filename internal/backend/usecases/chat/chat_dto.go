package chat

import "github.com/google/uuid"

type CreateChatInfo struct {
	Users []uuid.UUID `json:"users"`
}

type ChatInfo struct {
	Id    uuid.UUID   `json:"id"`
	Users []uuid.UUID `json:"users"`
}

type AddUserToChatInfo struct {
	UserId uuid.UUID `json:"nickname"`
	ChatId uuid.UUID `json:"chatId"`
}

type ChatInfoNoUsers struct {
	Id uuid.UUID `json:"id"`
}
