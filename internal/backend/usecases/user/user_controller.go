package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/helper"
)

type UserController interface {
	RegisterUser() http.Handler
	FindUsers() http.Handler
}

type userController struct {
	userRegistrationUseCase UserRegistrationUseCase
	findUsersUseCase        FindUsersUseCase
}

func NewUserController(userRegistrationUseCase UserRegistrationUseCase, findUsersUseCase FindUsersUseCase) UserController {
	return &userController{
		userRegistrationUseCase: userRegistrationUseCase,
		findUsersUseCase:        findUsersUseCase,
	}
}

func (controller *userController) RegisterUser() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userInfo, err := helper.Decode[UserRegistrationInfo](r)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = controller.userRegistrationUseCase.RegisterUser(userInfo)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			helper.EncodeNoBody(w, r, http.StatusOK)
		},
	)
}

func (controller *userController) FindUsers() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := r.PathValue("nickname")
			if len(nickname) == 0 {
				helper.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("nickname parameter shouldn't be empty"))
				return
			}
			users, err := controller.findUsersUseCase.FindUsers(nickname)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = helper.Encode(w, r, http.StatusOK, users)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}
