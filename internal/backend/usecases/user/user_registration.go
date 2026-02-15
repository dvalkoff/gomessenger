package user

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationUseCase interface {
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
		slog.Error("Failed to hash password", "error", err)
		return err
	}
	row := UserRow{
		Nickname:       userDto.Nickname,
		Name:           userDto.Name,
		HashedPassword: hashedPassword,
	}
	return useCase.userRepository.SaveUser(row)
}
