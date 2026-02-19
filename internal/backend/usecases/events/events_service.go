package events

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/usecases/chat"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSizeBytes = 50 * 1000 * 1000
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type EventsService interface {
	StreamEvents(StreamEventsDto) error
	CreateEvent(CreateEventDto) (EventCreatedDto, error)
}

type eventsService struct {
	chatRepository   chat.ChatRepository
	eventsRepository EventsRepository
	eventsHub        EventsHub
}

func NewEventService(chatRepository chat.ChatRepository, eventsRepository EventsRepository, eventsHub EventsHub) EventsService {
	return &eventsService{
		chatRepository:   chatRepository,
		eventsRepository: eventsRepository,
		eventsHub:        eventsHub,
	}
}

func (service *eventsService) CreateEvent(createEventDto CreateEventDto) (EventCreatedDto, error) {
	if len(createEventDto.Payload) > MaxMessageSizeBytes {
		slog.Error("Size of a message exceeds a limit", "message size", len(createEventDto.Payload))
		return EventCreatedDto{}, fmt.Errorf("Size of a message exceeds a limit")
	}
	userIds := []uuid.UUID{createEventDto.UserId}
	if createEventDto.ReceiverId != createEventDto.UserId {
		userIds = append(userIds, createEventDto.ReceiverId)
	}
	usersAreInChat, err := service.chatRepository.AreUsersInChat(createEventDto.ChatId, userIds)
	if err != nil {
		slog.Error("Failed to get chat participants", "error", err)
		return EventCreatedDto{}, err
	}
	if !usersAreInChat {
		slog.Error("Users are not participants of the chat", "userIds", userIds, "chatId", createEventDto.ChatId)
		return EventCreatedDto{}, fmt.Errorf("User is not participant of the chat")
	}
	eventRow := EventRow{
		UserId:     createEventDto.UserId,
		ReceiverId: createEventDto.ReceiverId,
		ChatId:     createEventDto.ChatId,
		Payload:    createEventDto.Payload,
	}
	eventId, err := service.eventsRepository.SaveEvent(eventRow)
	if err != nil {
		slog.Error("Failed to save event", "error", err)
		return EventCreatedDto{}, err
	}
	eventRow.Id = eventId

	service.eventsHub.SendEvent(eventRow)

	return EventCreatedDto{Id: eventId}, nil
}

func (service *eventsService) StreamEvents(dto StreamEventsDto) error {
	connection, err := wsUpgrader.Upgrade(dto.Writer, dto.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrate HTTP connection to Websockets", "error", err)
		return err
	}
	events, err := service.eventsRepository.GetEventsAfterCursor(dto.EventCursor, dto.UserId)
	if err != nil {
		slog.Error("Failed to get events for user", "userId", dto.UserId.String(), "error", err)
		return err
	}
	clientConnection := NewClientConnection(
		dto.UserId,
		connection,
	)
	go clientConnection.Run()
	for _, event := range events {
		clientConnection.SendEventToUser(event)
	}

	service.eventsHub.RegisterClient(clientConnection)

	return nil
}
