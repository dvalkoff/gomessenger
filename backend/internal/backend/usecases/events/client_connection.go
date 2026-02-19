package events

import (
	"log/slog"

	protoevent "github.com/dvalkoff/gomessenger/gen"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type ClientConnection interface {
	Run()
	SendEventToUser(EventRow)
	GetUserId() uuid.UUID
}

type clientConnection struct {
	userId     uuid.UUID
	connection *websocket.Conn
	events     chan EventRow
}

func NewClientConnection(userId uuid.UUID, connection *websocket.Conn) ClientConnection {
	return &clientConnection{
		userId:     userId,
		connection: connection,
		events:     make(chan EventRow),
	}
}

func (c *clientConnection) Run() {
	for {
		eventRow := <-c.events
		c.handleSendEventToUser(eventRow)
	}
}

func (c *clientConnection) SendEventToUser(event EventRow) {
	c.events <- event
}

func (c *clientConnection) handleSendEventToUser(event EventRow) {
	protoEvent := &protoevent.ServerEvent{
		Id:         &event.Id,
		UserId:     event.UserId[:],
		ReceiverId: event.ReceiverId[:],
		ChatId:     event.ChatId[:],
		Payload:    event.Payload,
	}
	bytes, err := proto.Marshal(protoEvent)
	if err != nil {
		slog.Error("Failed to marshall event", "eventId", event.Id, "error", err)
		return
	}
	err = c.connection.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		slog.Error("Failed to send event to websocket", "eventId", event.Id, "error", err)
	}
}
func (c *clientConnection) GetUserId() uuid.UUID {
	return c.userId
}
