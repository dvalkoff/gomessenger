package chat

import "log/slog"

type Chat struct {
	ChatId int
	messages chan Message
	activeClients map[*Client]struct{}
	register chan *Client
	unregister chan *Client
}

func NewChat(chatId int) *Chat {
	return &Chat{
		ChatId: chatId,
		messages: make(chan Message, 256),
		activeClients: map[*Client]struct{}{},
		register: make(chan *Client),
		unregister: make(chan *Client),
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
		case message := <-chat.messages:
			chat.sendMessageToClients(message)
		}
	}
}

func (chat *Chat) sendMessageToClients(message Message) {
	for client := range chat.activeClients {
		client.send <- message
	}
}