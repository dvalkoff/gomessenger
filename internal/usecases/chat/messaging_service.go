package chat

import (
	"log/slog"
	"net/http"
)

type ClientConnectionInfo struct {
	nickname string
	offset int
}

type MessagingService interface {
	CreateClient(cci ClientConnectionInfo, w http.ResponseWriter, r *http.Request) error
}

type messagingService struct {
	hub *Hub
	messagingRepository MessagingRepository
}

func NewMessagingService(hub *Hub, messagingRepository MessagingRepository) MessagingService {
	return &messagingService{
		hub: hub,
		messagingRepository: messagingRepository,
	}
}

func (service *messagingService) CreateClient(cci ClientConnectionInfo, w http.ResponseWriter, r *http.Request) error {
	messages, err := service.messagingRepository.GetMessages(cci.nickname, cci.offset)
	if err != nil {
		slog.Error("Failed to get messages", "error", err)
		return err
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrate HTTP connection to Websockets", "error", err)
		return err
	}
	client := &Client{
		nickname: cci.nickname,
		hub: service.hub,
		conn: conn,
		send: make(chan Message, 256),
	}

	go client.readMessages()
	go client.sendMessage()

	for _, messageRow := range messages {
		message := Message{
			Id: messageRow.id,
			MessageType: "message",
			ChatId: messageRow.chatId,
			Sender: messageRow.sender,
			Payload: messageRow.payload,
			SentAt: messageRow.sentAt,
		}
		client.send <- message
	}
	client.hub.registerClient <- client
	return nil
}