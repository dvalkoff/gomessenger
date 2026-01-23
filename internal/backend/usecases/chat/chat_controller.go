package chat

import (
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/utils"
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
			nickname := utils.GetNickname(r.Context())
			createChatInfo, err := utils.Decode[CreateChatInfo](r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			createChatInfo.CreatorNickname = nickname

			createdChat, err := controller.createChatUseCase.CreateChat(createChatInfo)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			if err = utils.Encode(w, r, http.StatusOK, createdChat); err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}

func (controller *chatController) GetChats() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := utils.GetNickname(r.Context())
			chats, err := controller.chatSelection.GetChats(nickname)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			err = utils.Encode(w, r, http.StatusOK, chats)
			if err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}

func (controller *chatController) AddUserToChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := utils.GetNickname(r.Context())
			addUserToChatInfo, err := utils.Decode[AddUserToChatInfo](r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			chatInfo, err := controller.createChatUseCase.AddUserToChat(nickname, addUserToChatInfo)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			if err := utils.Encode(w, r, http.StatusOK, chatInfo); err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
		},
	)
}
