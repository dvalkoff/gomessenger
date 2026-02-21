package chat

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type ChatService interface {
	CreateChat(CreateChatInfo, uuid.UUID) (ChatInfo, error)
	AddUserToChat(userId uuid.UUID, addUserToChatInfo AddUserToChatInfo) (ChatInfo, error)
	GetChats(userId uuid.UUID) ([]ChatInfoNoUsers, error)
}

type chatService struct {
	chatRepository ChatRepository
}

func NewChatService(chatRepository ChatRepository) ChatService {
	return &chatService{chatRepository: chatRepository}
}

func (service *chatService) CreateChat(createChatInfo CreateChatInfo, creator uuid.UUID) (ChatInfo, error) {
	chatId, err := uuid.NewRandom()
	if err != nil {
		slog.Error("Failed to generate uuid", "error", err)
		return ChatInfo{}, err
	}
	createChatInfo.Users = append(createChatInfo.Users, creator)
	chatRow := ChatRow{
		Id:    chatId,
		Users: createChatInfo.Users,
	}
	chat, err := service.chatRepository.CreateChat(chatRow)
	if err != nil {
		slog.Error("Failed to save chat", "error", err)
		return ChatInfo{}, err
	}
	return ChatInfo{
		Id:    chat.Id,
		Users: chat.Users,
	}, nil
}

func (service *chatService) AddUserToChat(userId uuid.UUID, addUserToChatInfo AddUserToChatInfo) (ChatInfo, error) {
	isUserInChat, err := service.chatRepository.AreUsersInChat(addUserToChatInfo.ChatId, []uuid.UUID{userId})
	if err != nil {
		slog.Error("Failed to get information about chat participants", "error", err)
		return ChatInfo{}, err
	}
	if !isUserInChat {
		slog.Error("User is not a chat participant")
		return ChatInfo{}, fmt.Errorf("User is not a chat participant")
	}
	chatUsers := ChatUsers{addUserToChatInfo.UserId}
	err = service.chatRepository.AddUsersToChat(addUserToChatInfo.ChatId, chatUsers)
	if err != nil {
		slog.Error("Failed to add user to chat", "error", err)
		return ChatInfo{}, err
	}
	chat, err := service.chatRepository.GetChat(addUserToChatInfo.ChatId)
	if err != nil {
		slog.Error("Failed to get chat", "error", err)
		return ChatInfo{}, err
	}
	return ChatInfo{
		Id:    chat.Id,
		Users: chat.Users,
	}, nil
}

func (service *chatService) GetChats(userId uuid.UUID) ([]ChatInfoNoUsers, error) {
	chatRows, err := service.chatRepository.GetChatsByUser(userId)
	if err != nil {
		slog.Error("Failed to get chats by userId", "userId", userId, "error", err)
		return nil, err
	}
	chats := make([]ChatInfoNoUsers, len(chatRows))
	for i, chatRow := range chatRows {
		chats[i] = ChatInfoNoUsers{
			Id: chatRow.Id,
		}
	}
	return chats, nil
}
