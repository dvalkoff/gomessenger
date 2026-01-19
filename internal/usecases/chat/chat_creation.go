package chat

import "fmt"

type CreateChatInfo struct {
	Name string `json:"name"`
	UserNickname string `json:"userNickname"`
}

type ChatInfo struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Users []ChatUser `json:"users"`
}

type ChatUser struct {
	Nickname string `json:"nickname"`
	Role string `json:"role"`
}

type AddUserToChatInfo struct {
	Nickname string `json:"nickname"`
	ChatId int `json:"chatId"`
}

type CreateChatUseCase interface {
	CreateChat(CreateChatInfo) (*ChatInfo, error)
	AddUserToChat(userNickname string, addUserToChatInfo AddUserToChatInfo) (*ChatInfo, error)
}

type createChatUseCase struct {
	chatRepository ChatRepository
}

func NewCreateChatUseCase(chatRepository ChatRepository) CreateChatUseCase {
	return &createChatUseCase{chatRepository: chatRepository}
}

func (useCase *createChatUseCase) CreateChat(createChatInfo CreateChatInfo) (*ChatInfo, error) {
	chat, err := useCase.chatRepository.CreateChat(createChatInfo.Name, createChatInfo.UserNickname)
	if err != nil {
		return nil, err
	}
	return mapChatInfo(chat), nil
}

func (useCase *createChatUseCase) AddUserToChat(user string, addUserToChatInfo AddUserToChatInfo) (*ChatInfo, error) {
	chat, err := useCase.chatRepository.GetChat(addUserToChatInfo.ChatId)
	if err != nil {
		return nil, err
	}
	userWhoAdds := findUserByNickname(chat.users, user)
	if userWhoAdds == nil || userWhoAdds.role != "admin" {
		return nil, fmt.Errorf("User is not administrator of a chat: %s", user)
	}
	chatUserRow, err := useCase.chatRepository.AddUserToChat(nil, addUserToChatInfo.ChatId, addUserToChatInfo.Nickname, "user")
	if err != nil {
		return nil, err
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

func mapChatInfo(chat *Chat) *ChatInfo {
	return &ChatInfo{
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
