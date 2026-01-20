package chat

import (
	"fmt"
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/helper"
)

type MessagingConrtoller interface {
	GetUpdates() http.Handler
	GetRealtimeUpdates() http.Handler
}

type messagingConrtoller struct {
	hub *Hub
}

func NewMessagingConrtoller(hub *Hub) MessagingConrtoller {
	return &messagingConrtoller{hub: hub}
}

func (controller *messagingConrtoller) GetUpdates() http.Handler {
	return nil
}

func (controller *messagingConrtoller) GetRealtimeUpdates() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := r.PathValue("nickname")
			if len(nickname) == 0 {
				helper.EncodeError(w, r, http.StatusBadRequest, fmt.Errorf("nickname parameter shouldn't be empty"))
				return
			}
			serve(nickname, controller.hub, w, r)
		},
	)
}
