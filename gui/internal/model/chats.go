package model

import "github.com/google/uuid"

type Chats struct {
	chats []ChatPreview
}

type ChatPreview struct {
	Id          uuid.UUID
	LastMessage string
	UnreadCount int
}
