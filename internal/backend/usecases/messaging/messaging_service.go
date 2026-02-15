package messaging

import (
	"log/slog"
	"net/http"
)

type MessagingService interface {
	CreateClient(cci ClientConnectionInfo, w http.ResponseWriter, r *http.Request) error
}

type messagingService struct {
	messagingHub        MessagingHub
	messagingRepository MessagingRepository
}

func NewMessagingService(messaingHub MessagingHub, messagingRepository MessagingRepository) MessagingService {
	return &messagingService{
		messagingHub:        messaingHub,
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
	client := &MessagingClient{
		nickname:     cci.nickname,
		messagingHub: service.messagingHub,
		conn:         conn,
		send:         make(chan Message, 256),
	}

	go client.readMessages()
	go client.sendMessage()

	for _, messageRow := range messages {
		message := Message{
			Id:        messageRow.id,
			EventType: "message",
			ChatId:    messageRow.chatId,
			Sender:    messageRow.sender,
			Payload:   messageRow.payload,
			SentAt:    messageRow.sentAt,
		}
		client.send <- message
	}
	client.messagingHub.RegisterClient(client)
	return nil
}
