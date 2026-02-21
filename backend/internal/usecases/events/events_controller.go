package events

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	protoevent "github.com/dvalkoff/gomessenger/backend/gen"
	"github.com/dvalkoff/gomessenger/backend/internal/utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	offsetQueryKey = "offset"
)

type EventsController interface {
	HandleEventsWS() http.Handler
	CreateEvent() http.Handler
}

type eventsController struct {
	eventsService EventsService
}

func NewEventsController(eventsService EventsService) EventsController {
	return &eventsController{eventsService: eventsService}
}

func (controller *eventsController) HandleEventsWS() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := utils.GetUserId(r.Context())
			cursor, err := getEventsCursor(r)
			if err != nil {
				EncodeProtoError(w, r, http.StatusBadRequest, err)
				return
			}
			dto := StreamEventsDto{
				UserId:      userId,
				EventCursor: cursor,
				Writer:      w,
				Request:     r,
			}
			err = controller.eventsService.StreamEvents(dto)
			if err != nil {
				EncodeProtoError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}

func getEventsCursor(r *http.Request) (int, error) {
	cursor := r.URL.Query().Get(offsetQueryKey)
	if len(cursor) == 0 {
		return 0, fmt.Errorf("Events cursor is not provided")
	}
	return strconv.Atoi(cursor)
}

func (controller *eventsController) CreateEvent() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			createEventDto, err := DecodeProtoEvent(r)
			if err != nil {
				slog.Error("Failed to decode proto message", "error", err)
				EncodeProtoError(w, r, http.StatusBadRequest, err)
				return
			}
			createdEventDto, err := controller.eventsService.CreateEvent(createEventDto)
			if err != nil {
				EncodeProtoError(w, r, http.StatusInternalServerError, err)
				return
			}
			EncodeProtoResponse(w, r, http.StatusOK, createdEventDto)
		},
	)
}

func DecodeProtoEvent(r *http.Request) (CreateEventDto, error) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return CreateEventDto{}, err
	}
	event := protoevent.ServerEvent{}
	err = proto.Unmarshal(bytes, &event)
	if err != nil {
		return CreateEventDto{}, err
	}
	userId, err := uuid.ParseBytes(event.UserId)
	if err != nil {
		return CreateEventDto{}, err
	}
	receiverId, err := uuid.ParseBytes(event.ReceiverId)
	if err != nil {
		return CreateEventDto{}, err
	}
	chatId, err := uuid.ParseBytes(event.ChatId)
	if err != nil {
		return CreateEventDto{}, err
	}
	return CreateEventDto{
		UserId:     userId,
		ReceiverId: receiverId,
		ChatId:     chatId,
		Payload:    event.Payload,
	}, nil
}

func EncodeProtoError(w http.ResponseWriter, r *http.Request, status int, errorToSend error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/protobuf")

	statusConverted := int32(status)
	errorConverted := errorToSend.Error()
	errorResponse := protoevent.ErrorResponse{
		Code:      &statusConverted,
		ErrorInfo: &errorConverted,
	}
	bytes, err := proto.Marshal(&errorResponse)
	if err != nil {
		slog.Error("Failed to marshall error response", "error", err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
}

func EncodeProtoResponse(w http.ResponseWriter, r *http.Request, status int, eventCreatedDto EventCreatedDto) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/protobuf")
	idConverted := uint64(eventCreatedDto.Id)
	response := protoevent.EventCreated{
		Id: &idConverted,
	}
	bytes, err := proto.Marshal(&response)
	if err != nil {
		slog.Error("Failed to marshall response", "error", err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
}
