package user

type UserRegistrationInfo struct {
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserInfo struct {
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
}
