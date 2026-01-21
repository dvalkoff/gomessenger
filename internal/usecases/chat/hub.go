package chat

import (
	"log/slog"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	chats map[int]*Chat
	clients map[string]*Client

	registerClient chan *Client
	unregisterClient chan *Client
	RegisterChat chan int
	UpdateClient chan string

	chatRepository ChatRepository
	messagingRepository MessagingRepository
}

func NewHub(chatRepository ChatRepository, messagingRepository MessagingRepository) *Hub {
	return &Hub{
		chats: make(map[int]*Chat),
		clients: make(map[string]*Client),
		registerClient: make(chan *Client),
		unregisterClient: make(chan *Client),
		RegisterChat: make(chan int),
		UpdateClient: make(chan string),
		chatRepository: chatRepository,
		messagingRepository: messagingRepository,
	}
}

func (h *Hub) Run() error {
	chatIds, err := h.chatRepository.GetChatIds()
	if err != nil {
		return err
	}
	for _, chatId := range chatIds {
		h.RegisterAndRunChat(chatId)
	}

	for {
		select {
		case nickname := <- h.UpdateClient:
			if client, ok := h.clients[nickname]; ok {
				h.RegisterClient(client)
			}
		case client := <-h.registerClient:
			h.RegisterClient(client)
		case client := <- h.unregisterClient:
			h.UnregisterClient(client)
		case chat := <- h.RegisterChat:
			h.RegisterAndRunChat(chat)
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	chatIds, err := h.chatRepository.GetChatIdsByUser(client.nickname)
	if err != nil {
		slog.Error("Failed to get chat ids by user", "nickname", client.nickname)
		return
	}
	for _, chatId := range chatIds {
		if chat, ok := h.chats[chatId]; ok {
			chat.register <- client
		}
	}
	h.clients[client.nickname] = client
}

func (h *Hub) UnregisterClient(client *Client) {
	chatIds, err := h.chatRepository.GetChatIdsByUser(client.nickname)
	if err != nil {
		slog.Error("Failed to get chat ids by user", "nickname", client.nickname)
		return
	}
	for _, chatid := range chatIds {
		h.chats[chatid].unregister <- client
	}
	close(client.send)
	delete(h.clients, client.nickname)
}

func (h *Hub) RegisterAndRunChat(chatId int) {
	chat := NewChat(chatId, h.messagingRepository)
	nicknames, err := h.chatRepository.GetNicknamesByChatId(chat.ChatId)
	if err != nil {
		slog.Error("Failed to get nicknames by chatId", "chatId", chat.ChatId)
		return
	}
	slog.Info("New chat has been registered", "chatId", chat.ChatId)
	h.chats[chat.ChatId] = chat
	go chat.Run()
	for _, nickname := range nicknames {
		if client, ok := h.clients[nickname]; ok {
			chat.register <- client
		}
	}
}
