package chat

import "log/slog"

type Chat struct {
	ChatId int
	messages chan Message
	activeClients map[*Client]struct{}
	register chan *Client
	unregister chan *Client

	messagingRepository MessagingRepository
}

func NewChat(chatId int, messagingRepository MessagingRepository) *Chat {
	return &Chat{
		ChatId: chatId,
		messages: make(chan Message, 256),
		activeClients: map[*Client]struct{}{},
		register: make(chan *Client),
		unregister: make(chan *Client),
		messagingRepository: messagingRepository,
	}
}

func (chat *Chat) RegisterClient(client *Client) {
	chat.register <- client
}

func (chat *Chat) UnregisterClient(client *Client) {
	chat.unregister <- client
}

func (chat *Chat) Run() {
	defer func() {
		close(chat.messages)
		close(chat.register)
		close(chat.unregister)
	}()
	for {
		select {
		case client := <-chat.register:
			chat.activeClients[client] = struct{}{}
			slog.Info("New client was registered to the chat", "chatId", chat.ChatId, "nickname", client.nickname)
		case client := <-chat.unregister:
			delete(chat.activeClients, client)
			slog.Info("Client was unregistered from the chat", "chatId", chat.ChatId, "nickname", client.nickname)
		case message := <-chat.messages:
			chat.sendMessageToClients(message)
		}
	}
}

func (chat *Chat) sendMessageToClients(message Message) {
	row, err := chat.messagingRepository.SaveMessage(message)
	if err != nil {
		slog.Error("Failed to save message", "error", err)
		return
	}
	message.Id = row.id
	for client := range chat.activeClients {
		select {
		case client.send <- message:
		default:
			slog.Warn("Failed to send message to a client", "chat", chat.ChatId, "nickname", client.nickname)
		}
	}
}
