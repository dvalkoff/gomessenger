package events

import (
	"log/slog"

	"github.com/google/uuid"
)

type EventsHub interface {
	Run()
	RegisterClient(client ClientConnection)
	// UnregisterClient(client *MessagingClient)
	SendEvent(event EventRow)
	// Shutdown() <-chan struct{}
}

type eventsHub struct {
	clientConnections map[uuid.UUID]ClientConnection
	events            chan EventRow
	registerClients   chan ClientConnection
}

func NewEventsHub() EventsHub {
	return &eventsHub{
		clientConnections: make(map[uuid.UUID]ClientConnection),
		events:            make(chan EventRow),
		registerClients:   make(chan ClientConnection),
	}
}

func (hub *eventsHub) Run() {
	for {
		select {
		case event := <-hub.events:
			hub.processEvent(event)
		case client := <-hub.registerClients:
			hub.processRegisterClient(client)
		}
	}
}

func (hub *eventsHub) SendEvent(event EventRow) {
	hub.events <- event
}

func (hub *eventsHub) processEvent(event EventRow) {
	clientConnection, ok := hub.clientConnections[event.ReceiverId]
	if !ok {
		slog.Debug("Client is offline. Event will be fetched later", "eventId", event.Id)
		return
	}
	clientConnection.SendEventToUser(event)
}

func (hub *eventsHub) RegisterClient(client ClientConnection) {
	hub.registerClients <- client
}

func (hub *eventsHub) processRegisterClient(client ClientConnection) {
	if existingClient, ok := hub.clientConnections[client.GetUserId()]; ok {
		existingClient.Run() // TODO: need to close existing connection
	}
	hub.clientConnections[client.GetUserId()] = client
}
