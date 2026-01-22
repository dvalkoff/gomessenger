package user

import (
	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationInfo struct {
	Nickname string `json:"nickname"`
	Name string `json:"name"`
	Password string `json:"password"`
}

type UserRegistrationUseCase interface{
	RegisterUser(UserRegistrationInfo) error
}

type userRegistrationUseCase struct {
	userRepository UserRepository
}

func NewUserUserRegistrationUseCase(userRepository UserRepository) UserRegistrationUseCase {
	return &userRegistrationUseCase{userRepository: userRepository}
}

func (useCase *userRegistrationUseCase) RegisterUser(userDto UserRegistrationInfo) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	row := UserRow{
		nickname: userDto.Nickname,
		name: userDto.Name,
		hashedPassword: hashedPassword,
	}
	return useCase.userRepository.saveUser(row)
}
