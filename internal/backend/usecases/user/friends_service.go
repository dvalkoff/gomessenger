package user

import "context"

type FriendsService interface {
	AddFriend(ctx context.Context, userNickname string, friendNickname string) error
	GetFriends(ctx context.Context, userNickname string) ([]UserInfo, error)
}

type friendsService struct {
	userRepository UserRepository
}

func NewFriendsService(userRepository UserRepository) FriendsService {
	return &friendsService{userRepository: userRepository}
}

func (service *friendsService) AddFriend(ctx context.Context, userNickname string, friendNickname string) error {
	return service.userRepository.AddFriend(ctx, userNickname, friendNickname)
}

func (service *friendsService) GetFriends(ctx context.Context, userNickname string) ([]UserInfo, error) {
	userRows, err := service.userRepository.GetFriends(ctx, userNickname)
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
