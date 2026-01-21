package chat

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Id int `json:"id"`
	MessageType string `json:"messageType"`
	ChatId int `json:"chatId"`
	Sender string `json:"sender"`
	Payload string `json:"payload"`
	SentAt time.Time `json:"sentAt"`
}

type Client struct {
	nickname string
	hub *Hub
	conn *websocket.Conn
	send chan Message
}

func (c *Client) sendMessage() {
	defer func() {
		c.hub.unregisterClient <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, reader, err := c.conn.NextReader()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Info("Connection was closed by a client", "error", err)
			} else {
				slog.Error("Failed to read from socket", "error", err)
			}
			return
		}

		message := Message{}
		err = json.NewDecoder(reader).Decode(&message)
		if err != nil {
			slog.Error("Failed to decode message to json", "error", err)
			return
		}
		message.Sender = c.nickname
		message.SentAt = time.Now()
		message.Payload = strings.TrimSpace(message.Payload)
		c.hub.chats[message.ChatId].messages <- message // TODO: remove map access, replace with channels
	}
}

func (c *Client) readMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				slog.Info("Messages channel was closed for client", "nickname", c.nickname)
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteJSON(message)
			if err != nil {
				slog.Error("Failed to write message to connection", "error", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Failed to ping client", "error", err)
				return
			}
		}
	}
}
