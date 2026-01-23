package messaging

import (
	"net/http"

	"github.com/dvalkoff/gomessenger/internal/backend/utils"
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
			nickname := utils.GetNickname(r.Context())
			cci := ClientConnectionInfo{
				nickname: nickname,
				offset:   0,
			}
			err := controller.messagingService.CreateClient(cci, w, r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}
