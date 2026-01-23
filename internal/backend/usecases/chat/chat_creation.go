package chat

import (
	"fmt"
)

type CreateChatUseCase interface {
	CreateChat(CreateChatInfo) (ChatInfo, error)
	AddUserToChat(userNickname string, addUserToChatInfo AddUserToChatInfo) (ChatInfo, error)
}

type createChatUseCase struct {
	chatRepository ChatRepository
}

func NewCreateChatUseCase(chatRepository ChatRepository) CreateChatUseCase {
	return &createChatUseCase{chatRepository: chatRepository}
}

func (useCase *createChatUseCase) CreateChat(createChatInfo CreateChatInfo) (ChatInfo, error) {
	chatUsers := make([]ChatUserRow, 0)
	chatUsers = append(chatUsers, ChatUserRow{createChatInfo.CreatorNickname, "admin"})
	for _, userToAdd := range createChatInfo.Users {
		chatUsers = append(chatUsers, ChatUserRow{userToAdd, "user"})
	}
	chatRow := ChatRow{
		name: createChatInfo.Name,
		users: chatUsers,
	}
	chat, err := useCase.chatRepository.CreateChat(chatRow)
	if err != nil {
		return ChatInfo{}, err
	}
	return mapChatInfo(chat), nil
}

func (useCase *createChatUseCase) AddUserToChat(user string, addUserToChatInfo AddUserToChatInfo) (ChatInfo, error) {
	chat, err := useCase.chatRepository.GetChat(addUserToChatInfo.ChatId)
	if err != nil {
		return ChatInfo{}, err
	}
	userWhoAdds := findUserByNickname(chat.users, user)
	if userWhoAdds == nil || userWhoAdds.role != "admin" {
		return ChatInfo{}, fmt.Errorf("User is not administrator of a chat: %s", user)
	}
	chatUserRow := ChatUserRow{
		nickname: addUserToChatInfo.Nickname,
		role: "user",
	}
	err = useCase.chatRepository.AddUsersToChat(addUserToChatInfo.ChatId, []ChatUserRow{chatUserRow})
	if err != nil {
		return ChatInfo{}, err
	}
	chat.users = append(chat.users, chatUserRow)
	return mapChatInfo(chat), nil
}

func findUserByNickname(users []ChatUserRow, nickname string) *ChatUserRow {
	for _, user := range users {
		if user.nickname == nickname {
			return &user
		}
	}
	return nil
}

func mapChatInfo(chat ChatRow) ChatInfo {
	return ChatInfo{
		Id: chat.id,
		Name: chat.name,
		Users: mapUsers(chat.users),
	}
}

func mapUsers(users []ChatUserRow) []ChatUser {
	mappedUsers := make([]ChatUser, len(users))
	for i, user := range users {
		mappedUsers[i] = ChatUser{
			Nickname: user.nickname,
			Role: user.role,
		}
	}
	return mappedUsers
}
