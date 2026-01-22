package chat

import "log/slog"

type ChatInfoNoUsers struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type ChatSelection interface {
	GetChats(nickname string) ([]ChatInfoNoUsers, error)
}

type chatSelection struct {
	chatRepository ChatRepository
}

func NewChatSelection(chatRepository ChatRepository) ChatSelection {
	return &chatSelection{chatRepository: chatRepository}
}

func (cs *chatSelection) GetChats(nickname string) ([]ChatInfoNoUsers, error) {
	chatRows, err := cs.chatRepository.GetChatsNoUsersByNickname(nickname)
	if err != nil {
		slog.Error("Failed to get chats by nickname", "nickname", nickname, "error", err)
		return nil, err
	}
	chats := make([]ChatInfoNoUsers, len(chatRows))
	for i, chatRow := range chatRows {
		chats[i] = ChatInfoNoUsers{
			Id: chatRow.id,
			Name: chatRow.name,
		}
	}
	return chats, nil
}
