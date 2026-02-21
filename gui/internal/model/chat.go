package model

import "github.com/google/uuid"

type Chat struct {
	Id           uuid.UUID
	Participants []ChatParticipant
	Messages     []ChatMessage
}

type ChatParticipant struct {
	Id       uuid.UUID
	Nickname string
}

type ChatMessage struct {
	Id       int
	SenderId uuid.UUID
	Payload  string
}
