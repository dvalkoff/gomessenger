package chat

import (
	"log/slog"
	"net/http"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	MessageType string `json:"messageType"`
	ChatId int `json:"chatId"`
	Sender string `json:"sender"`
	Payload string `json:"payload"`
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
		message := Message{}
		message.Sender = c.nickname
		err := c.conn.ReadJSON(&message)
		if err != nil {
			slog.Error("Failed to send message", "error", err)
		}

		message.Payload = strings.TrimSpace(message.Payload)
		c.hub.chats[message.ChatId].messages <- message
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
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				slog.Error("Failed to write message to connection", "error", err)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serve(nickname string, hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrate HTTP connection to Websockets", "error", err)
		return
	}
	client := &Client{
		nickname: nickname,
		hub: hub,
		conn: conn,
		send: make(chan Message, 256),
	}
	client.hub.registerClient <- client

	go client.readMessages()
	go client.sendMessage()
}
