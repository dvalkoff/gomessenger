package user

type FriendsService interface {
	AddFriend(userNickname string, friendNickname string) error
	GetFriends(userNickname string) ([]UserInfo, error)
}

type friendsService struct {
	userRepository UserRepository
}

func NewFriendsService(userRepository UserRepository) FriendsService {
	return &friendsService{userRepository: userRepository}
}

func (service *friendsService) AddFriend(userNickname string, friendNickname string) error {
	return service.userRepository.AddFriend(userNickname, friendNickname)
}

func (service *friendsService) GetFriends(userNickname string) ([]UserInfo, error) {
	userRows, err := service.userRepository.GetFriends(userNickname)
	if err != nil {
		return nil, err
	}
	users := make([]UserInfo, len(userRows))
	for i, userRow := range userRows {
		users[i] = UserInfo{
			Nickname: userRow.Nickname,
			Name:     userRow.Name,
		}
	}
	return users, nil
}
