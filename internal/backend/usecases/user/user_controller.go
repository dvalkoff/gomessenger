package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/utils"
	"github.com/google/uuid"
)

type UserController interface {
	RegisterUser() http.Handler
	FindUsers() http.Handler
	AddContact() http.Handler
	GetContacts() http.Handler
}

type userController struct {
	userService     UserService
	contactsService ContactsService
}

func NewUserController(userService UserService, contactsServer ContactsService) UserController {
	return &userController{
		userService:     userService,
		contactsService: contactsServer,
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
			responseDto, err := controller.userService.RegisterUser(userInfo)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = utils.Encode(w, r, http.StatusOK, responseDto)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
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
			users, err := controller.userService.FindUsers(nickname)
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

func (controller *userController) AddContact() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := utils.GetUserId(r.Context())
			contactUserId, err := uuid.Parse(r.PathValue("contactId"))
			if err != nil {
				utils.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("contactId parameter is invalid: %w", err))
				return
			}
			err = controller.contactsService.AddContact(userId, contactUserId)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			utils.EncodeNoBody(w, r, 200)
		},
	)
}

func (controller *userController) GetContacts() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := utils.GetUserId(r.Context())
			users, err := controller.contactsService.GetContacts(userId)
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
