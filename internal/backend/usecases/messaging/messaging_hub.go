package messaging

import (
	"log/slog"
	"slices"

	"github.com/dvalkoff/gomessenger/internal/backend/usecases/chat"
)

type MessagingHub interface {
	Run()
	RegisterClient(client *MessagingClient)
	UnregisterClient(client *MessagingClient)
	SendMessage(message Message)
	Shutdown()
}

type messagingHub struct {
	clients map[string][]*MessagingClient

	registerClientChan   chan *MessagingClient
	unregisterClientChan chan *MessagingClient

	messagesChan chan Message
	shutdownChan chan struct{}

	chatRepository      chat.ChatRepository
	messagingRepository MessagingRepository
}

func NewMessagingHub(chatRepository chat.ChatRepository, messagingRepository MessagingRepository) MessagingHub {
	return &messagingHub{
		clients: make(map[string][]*MessagingClient),

		registerClientChan:   make(chan *MessagingClient),
		unregisterClientChan: make(chan *MessagingClient),

		messagesChan: make(chan Message),
		shutdownChan: make(chan struct{}),

		chatRepository:      chatRepository,
		messagingRepository: messagingRepository,
	}
}

func (h *messagingHub) Run() {
	for {
		select {
		case client := <-h.registerClientChan:
			h.processRegisterClient(client)
		case client := <-h.unregisterClientChan:
			h.processUnregisterClient(client)
		case message := <-h.messagesChan:
			h.processSendMessage(message)
		case <-h.shutdownChan:
			h.processShutdown()
		}
	}
}

func (h *messagingHub) RegisterClient(client *MessagingClient) {
	h.registerClientChan <- client
}

func (h *messagingHub) UnregisterClient(client *MessagingClient) {
	h.unregisterClientChan <- client
}

func (h *messagingHub) SendMessage(message Message) {
	h.messagesChan <- message
}

func (h *messagingHub) Shutdown() {
	h.shutdownChan <- struct{}{}
}

func (h *messagingHub) processRegisterClient(client *MessagingClient) {
	h.clients[client.nickname] = append(h.clients[client.nickname], client)
}

func (h *messagingHub) processUnregisterClient(client *MessagingClient) {
	if len(h.clients[client.nickname]) <= 1 {
		delete(h.clients, client.nickname)
	} else {
		oldSlice := h.clients[client.nickname]
		clientIdx := 0
		for i, cl := range oldSlice {
			if cl == client {
				clientIdx = i
			}
		}
		newSlice := make([]*MessagingClient, 0)
		newSlice = append(newSlice, oldSlice[:clientIdx]...)
		h.clients[client.nickname] = append(newSlice, oldSlice[clientIdx+1:]...)
	}
	close(client.send)
}

func (h *messagingHub) processSendMessage(message Message) {
	nicknames, err := h.chatRepository.GetNicknamesByChatId(message.ChatId)
	if err != nil {
		slog.Error("Failed to get chat users", "error", err)
		return
	}
	if !slices.Contains(nicknames, message.Sender) {
		slog.Warn("User is not allowed to send a message to the chat", "chatId", message.ChatId, "nickname", message.Sender)
		return
	}

	row, err := h.messagingRepository.SaveMessage(MessageRow{
		payload: message.Payload,
		sender: message.Sender,
		chatId: message.ChatId,
		sentAt: message.SentAt,
	})
	if err != nil {
		slog.Error("Failed to save message", "error", err)
		return
	}
	message.Id = row.id

	for _, nickname := range nicknames {
		if clients, ok := h.clients[nickname]; ok {
			for _, client := range clients {
				select {
				case client.send <- message:
				default:
					slog.Warn("Failed to send message to a client", "chat", message.ChatId, "nickname", client.nickname)
				}
			}
		}
	}
}

func (h *messagingHub) processShutdown() {
	for _, clients := range h.clients {
		for _, client := range clients {
			h.processUnregisterClient(client)
		}
	}
	slog.Info("Messaging hub was successfuly shut down")
}
