package messaging

import (
	"database/sql"
	"time"
)

type MessageRow struct {
	id      int
	payload string
	sender  string
	chatId  int
	sentAt  time.Time
}

type MessagingRepository interface {
	GetMessages(nickname string, offset int) ([]MessageRow, error)
	SaveMessage(message MessageRow) (MessageRow, error)
}

type messagingRepository struct {
	db *sql.DB
}

func NewMessagingRepository(db *sql.DB) MessagingRepository {
	return &messagingRepository{db: db}
}

func (repository *messagingRepository) SaveMessage(message MessageRow) (MessageRow, error) {
	sql := `INSERT INTO messenger.messages(payload, sender, chat_id, sent_at)
	VALUES($1, $2, $3, $4) RETURNING id`
	err := repository.db.
		QueryRow(sql, message.payload, message.sender, message.chatId, message.sentAt).
		Scan(&message.id)
	if err != nil {
		return MessageRow{}, err
	}
	return message, nil
}

func (repository *messagingRepository) GetMessages(nickname string, offset int) ([]MessageRow, error) {
	sql := `
	SELECT m.id, m.payload, m.sender, m.chat_id, m.sent_at FROM messenger.messages m
	JOIN messenger.chats_users cu ON cu.chat_id = m.chat_id
	WHERE cu.user_nickname = $1 AND m.id > $2
	ORDER BY m.id`
	rows, err := repository.db.Query(sql, nickname, offset)
	if err != nil {
		return nil, err
	}
	messageRows := make([]MessageRow, 0)
	for rows.Next() {
		messageRow := MessageRow{}
		err = rows.Scan(&messageRow.id, &messageRow.payload, &messageRow.sender, &messageRow.chatId, &messageRow.sentAt)
		if err != nil {
			return nil, err
		}
		messageRows = append(messageRows, messageRow)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messageRows, nil
}
