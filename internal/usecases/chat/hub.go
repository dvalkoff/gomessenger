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
	RegisterChat chan *Chat

	chatRepository ChatRepository
}

func NewHub(chatRepository ChatRepository) *Hub {
	return &Hub{
		chats: make(map[int]*Chat),
		clients: make(map[string]*Client),
		registerClient: make(chan *Client),
		unregisterClient: make(chan *Client),
		RegisterChat: make(chan *Chat),
		chatRepository: chatRepository,
	}
}

func (h *Hub) Run() error {
	chatIds, err := h.chatRepository.GetChatIds()
	if err != nil {
		return err
	}
	for _, chatId := range chatIds {
		h.RegisterAndRunChat(NewChat(chatId))
	}

	for {
		select {
		case client := <-h.registerClient:
			h.RegisterClient(client)
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

func (h *Hub) RegisterAndRunChat(chat *Chat) {
	nicknames, err := h.chatRepository.GetNicknamesByChatId(chat.ChatId)
	if err != nil {
		slog.Error("Failed to get nicknames by chatId", "chatId", chat.ChatId)
		return
	}
	for _, nickname := range nicknames {
		if client, ok := h.clients[nickname]; ok {
			chat.activeClients[client] = struct{}{}
		}
	}
	slog.Info("New chat has been registered", "chatId", chat.ChatId)
	h.chats[chat.ChatId] = chat
	go chat.Run()
}
