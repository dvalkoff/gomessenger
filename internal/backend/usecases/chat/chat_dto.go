package chat

type CreateChatInfo struct {
	Name            string   `json:"name"`
	CreatorNickname string   `json:"creatorNickname"`
	Users           []string `json:"users"`
}

type ChatInfo struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	Users []ChatUser `json:"users"`
}

type ChatUser struct {
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type AddUserToChatInfo struct {
	Nickname string `json:"nickname"`
	ChatId   int    `json:"chatId"`
}

type ChatInfoNoUsers struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
