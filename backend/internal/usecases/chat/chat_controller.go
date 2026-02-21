package chat

import (
	"log/slog"
	"net/http"

	"github.com/dvalkoff/gomessenger/backend/internal/utils"
)

type ChatController interface {
	CreateChat() http.Handler
	GetChats() http.Handler
	AddUserToChat() http.Handler
}

type chatController struct {
	chatService ChatService
}

func NewChatController(createChatUseCase ChatService) ChatController {
	return &chatController{chatService: createChatUseCase}
}

func (controller *chatController) CreateChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := utils.GetUserId(r.Context())
			createChatInfo, err := utils.Decode[CreateChatInfo](r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusBadRequest, err)
				return
			}
			createdChat, err := controller.chatService.CreateChat(createChatInfo, userId)
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
			userId := utils.GetUserId(r.Context())
			chats, err := controller.chatService.GetChats(userId)
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
			userId := utils.GetUserId(r.Context())
			addUserToChatInfo, err := utils.Decode[AddUserToChatInfo](r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
			chatInfo, err := controller.chatService.AddUserToChat(userId, addUserToChatInfo)
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
