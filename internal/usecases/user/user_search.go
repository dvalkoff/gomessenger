package user

type UserInfo struct {
	Nickname string `json:"nickname"`
	Name string `json:"name"`
}

type FindUsersUseCase interface {
	FindUsers(string) ([]UserInfo, error)
}

type findUsersUseCase struct {
	userRepository UserRepository
}

func NewFindUsersUseCase(userRepository UserRepository) FindUsersUseCase {
	return &findUsersUseCase{userRepository: userRepository}
}

func (useCase *findUsersUseCase) FindUsers(nicknameSubstring string) ([]UserInfo, error) {
	users, err := useCase.userRepository.findUsersByNickname(nicknameSubstring)
	if err != nil {
		return nil, err
	}
	// TODO: rewrite to map
	mappedUsers := make([]UserInfo, len(users))
	for i, user := range users {
		mappedUsers[i] = UserInfo{
			Nickname: user.nickname,
			Name: user.name,
		}
	}
	return mappedUsers, nil
}
