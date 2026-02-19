package events

import (
	"net/http"

	"github.com/google/uuid"
)

type CreateEventDto struct {
	UserId     uuid.UUID
	ReceiverId uuid.UUID
	ChatId     uuid.UUID
	Payload    []byte
}

type EventCreatedDto struct {
	Id uint64
}

type StreamEventsDto struct {
	UserId      uuid.UUID
	EventCursor int
	Writer      http.ResponseWriter
	Request     *http.Request
}
