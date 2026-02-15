package user

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
	users, err := useCase.userRepository.FindUsersByNickname(nicknameSubstring)
	if err != nil {
		return nil, err
	}
	mappedUsers := make([]UserInfo, len(users))
	for i, user := range users {
		mappedUsers[i] = UserInfo{
			Nickname: user.Nickname,
			Name:     user.Name,
		}
	}
	return mappedUsers, nil
}
