package chat

import (
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/helper"
)

type ChatController interface {
	CreateChat() http.Handler
	GetChats() http.Handler
	AddUserToChat() http.Handler
}

type chatController struct {
	createChatUseCase CreateChatUseCase
	chatSelection     ChatSelection
}

func NewChatController(createChatUseCase CreateChatUseCase, chatSelection ChatSelection) ChatController {
	return &chatController{createChatUseCase: createChatUseCase, chatSelection: chatSelection}
}

func (controller *chatController) CreateChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := helper.GetNickname(r.Context())
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

func (controller *chatController) GetChats() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := helper.GetNickname(r.Context())
			chats, err := controller.chatSelection.GetChats(nickname)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = helper.Encode(w, r, http.StatusOK, chats)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}

func (controller *chatController) AddUserToChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := helper.GetNickname(r.Context())
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
