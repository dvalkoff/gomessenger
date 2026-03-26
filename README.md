# Go Messenger Backend

A real-time messaging backend built with pure Go. Supports user authentication, friend management, group chats, and live messaging over WebSockets.

---

## Features

- **User registration & login** with secure password hashing
- **JWT-based authentication** for stateless, secure API access
- **Friend system** — send and manage friend connections
- **Group chats** — create conversations with multiple participants
- **Real-time messaging** via WebSockets

---

## Tech Stack

| Layer | Technology |
|---|---|
| HTTP Server | Pure Go (`net/http`) |
| WebSockets | [Gorilla WebSocket](https://github.com/gorilla/websocket) |
| Database Driver | [pq](https://github.com/lib/pq) |
| Database | PostgreSQL |
| Migrations | [Liquibase](https://www.liquibase.org/) |

---

## Getting Started

```bash
docker compose up -d
```

This will spin up three services in order:

1. **postgres** — PostgreSQL database
2. **liquibase** — Liquibase runs all pending migrations, then exits
3. **messenger** — the Go server, started after migrations complete

The API will be available at `http://localhost:8080`.

---

## API Reference

All protected endpoints require a `Bearer` token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

### Auth

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/signup` | ❌ | Register a new user |
| `POST` | `/signin` | ❌ | Log in and receive a JWT |

#### Register — `POST /signup`

```json
{
  "nickname": "alice",
  "name": "alice",
  "password": "password"
}
```

#### Login — `POST /signin`

```json
{
  "nickname": "alice",
  "password": "password"
}
```

**Response:**
```json
{
  "accessToken": "<jwt token>"
}
```

---

### Friends

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/users/friends` | ✅ | Add a friend |
| `GET` | `/users/friends` | ✅ | List all friends |

---

### Chats

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/chats` | ✅ | Create a new chat |
| `GET` | `/chats` | ✅ | List all chats for the current user |

#### Create Chat — `POST /api/chats`

```json
{
  "name": "Chat name",
  "users": ["nickname1", "nickname2"]
}
```

---

### Real-time Messaging (WebSocket)

Connect to the WebSocket endpoint to send and receive messages in real time:

```
ws://localhost:8080/messaging?token=<your_jwt_token>
```

Authentication is performed via the `token` query parameter on the initial handshake. The connection is rejected if the token is missing or invalid.

#### Sending a message

```json
{
  "chat_id": 1,
  "payload": "Hey everyone!"
}
```

#### Receiving a message

```json
{
  "id": 1,
  "chat_id": 1,
  "sender": "sender_nickname",
  "payload": "Hey everyone!",
  "sentAt": "2025-03-26T14:00:00Z"
}
```

---

## Project Structure

```
go-messenger/
├── cmd/
│   └── api/             # Application entrypoint
│       └── main.go
├── internal/
│   ├── config           # Application configuration (DB, HTTP server)
│   ├── middleware/      # Auth middleware, CORS, Logging
│   ├── utils/           # Common utils package
│   └── usecases/        # Main functionality
│       ├── chat/        # Chats management
│       ├── messaging/   # WebSocket hub and connection management
│       └── user/        # User-related business logic
├── migrations/
│   └── changelog/       # Liquibase migration files
│       └── db.changelog-master.xml
├── .env
├── compose.yaml
├── go.mod
├── go.sum
└── README.md
```
