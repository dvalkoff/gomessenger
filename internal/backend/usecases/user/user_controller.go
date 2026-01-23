package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/utils"
)

type UserController interface {
	RegisterUser() http.Handler
	FindUsers() http.Handler
	AddFriend() http.Handler
	GetFriends() http.Handler
}

type userController struct {
	userRegistrationUseCase UserRegistrationUseCase
	findUsersUseCase        FindUsersUseCase
	friendsService FriendsService
}

func NewUserController(userRegistrationUseCase UserRegistrationUseCase, findUsersUseCase FindUsersUseCase, friendsService FriendsService) UserController {
	return &userController{
		userRegistrationUseCase: userRegistrationUseCase,
		findUsersUseCase:        findUsersUseCase,
		friendsService: friendsService,
	}
}

func (controller *userController) RegisterUser() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userInfo, err := utils.Decode[UserRegistrationInfo](r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = controller.userRegistrationUseCase.RegisterUser(userInfo)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			utils.EncodeNoBody(w, r, http.StatusOK)
		},
	)
}

func (controller *userController) FindUsers() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := r.PathValue("nickname")
			if len(nickname) == 0 {
				utils.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("nickname parameter shouldn't be empty"))
				return
			}
			users, err := controller.findUsersUseCase.FindUsers(nickname)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = utils.Encode(w, r, http.StatusOK, users)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}

func (controller *userController) AddFriend() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := utils.GetNickname(r.Context())
			friendsNickname := r.PathValue("friendsNickname")
			if len(friendsNickname) == 0 {
				utils.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("friendsNickname parameter shouldn't be empty"))
				return
			}
			err := controller.friendsService.AddFriend(nickname, friendsNickname)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			utils.EncodeNoBody(w, r, 200)
		},
	)
}

func (controller *userController) GetFriends() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := utils.GetNickname(r.Context())
			users, err := controller.friendsService.GetFriends(nickname)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
			}
			err = utils.Encode(w, r, http.StatusOK, users)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}
