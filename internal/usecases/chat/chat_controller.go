package chat

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/helper"
)

type ChatController interface {
	CreateChat() http.Handler
	AddUserToChat() http.Handler
}

type chatController struct {
	createChatUseCase CreateChatUseCase
}

func NewChatController(createChatUseCase CreateChatUseCase) ChatController {
	return &chatController{createChatUseCase: createChatUseCase}
}

func (controller *chatController) CreateChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := r.PathValue("nickname")
			if len(nickname) == 0 {
				helper.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("nickname parameter shouldn't be empty"))
				return
			}
			createChatInfo, err := helper.Decode[CreateChatInfo](r)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			createChatInfo.CreatorNickname = nickname

			createdChat, err := controller.createChatUseCase.CreateChat(createChatInfo)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			if err = helper.Encode(w, r, http.StatusOK, createdChat); err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}

func (controller *chatController) AddUserToChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := r.PathValue("nickname")
			if len(nickname) == 0 {
				helper.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("nickname parameter shouldn't be empty"))
				return
			}
			addUserToChatInfo, err := helper.Decode[AddUserToChatInfo](r)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			chatInfo, err := controller.createChatUseCase.AddUserToChat(nickname, addUserToChatInfo)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			if err := helper.Encode(w, r, http.StatusOK, chatInfo); err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}