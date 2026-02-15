package messaging

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dvalkoff/gomessenger/internal/backend/utils"
)

const (
	offsetQueryKey = "offset"
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
			offset, err := getMessagesOffset(r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusBadRequest, err)
				return
			}
			cci := ClientConnectionInfo{
				nickname: nickname,
				offset:   offset,
			}
			err = controller.messagingService.CreateClient(cci, w, r)
			if err != nil {
				utils.EncodeError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}

func getMessagesOffset(r *http.Request) (int, error) {
	offset := r.URL.Query().Get(offsetQueryKey)
	if len(offset) == 0 {
		return 0, fmt.Errorf("Message offset is not provided")
	}
	return strconv.Atoi(offset)
}
