package messaging

import (
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/helper"
)

type MessagingConrtoller interface {
	GetRealtimeUpdates() http.Handler
}

type messagingConrtoller struct {
	messagingService MessagingService
}

func NewMessagingConrtoller(messagingService MessagingService) MessagingConrtoller {
	return &messagingConrtoller{messagingService: messagingService}
}

func (controller *messagingConrtoller) GetRealtimeUpdates() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname := helper.GetNickname(r.Context())
			cci := ClientConnectionInfo{
				nickname: nickname,
				offset:   0,
			}
			err := controller.messagingService.CreateClient(cci, w, r)
			if err != nil {
				helper.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}
