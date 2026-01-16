package user

import (
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/helper"
)

type UserController interface {
	RegisterUser() http.Handler
}

type userController struct {
	userRegistrationUseCase UserRegistrationUseCase
}

func NewUserController(userRegistrationUseCase UserRegistrationUseCase) UserController {
	return &userController{userRegistrationUseCase: userRegistrationUseCase}
}

func (controller *userController) RegisterUser() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userInfo, err := helper.Decode[UserRegistrationInfo](r)
			if err != nil {
				helper.EncodeError(w, r, 500, err)
				return
			}
			err = controller.userRegistrationUseCase.RegisterUser(userInfo)
			if err != nil {
				helper.EncodeError(w, r, 500, err)
				return
			}
			helper.EncodeNoBody(w, r, 200)
		},
	)
}