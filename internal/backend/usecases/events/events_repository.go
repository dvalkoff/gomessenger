package events

import (
	"database/sql"

	"github.com/google/uuid"
)

type EventRow struct {
	Id         uint64
	UserId     uuid.UUID
	ReceiverId uuid.UUID
	ChatId     uuid.UUID
	Payload    []byte
}

type EventsRepository interface {
	SaveEvent(EventRow) (uint64, error)
	GetEventsAfterCursor(cursor int, receiverId uuid.UUID) ([]EventRow, error)
}

type eventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) EventsRepository {
	return &eventsRepository{db: db}
}

func (repository *eventsRepository) SaveEvent(row EventRow) (uint64, error) {
	sql := `INSERT INTO messenger.events(user_id, receiver_id, chat_id, payload)
			VALUES($1, $2, $3, $4) RETURNING id`
	var eventId uint64
	err := repository.db.
		QueryRow(sql, row.UserId, row.ReceiverId, row.ChatId, row.Payload).
		Scan(&eventId)
	if err != nil {
		return 0, err
	}
	return eventId, nil
}

func (repository *eventsRepository) GetEventsAfterCursor(cursor int, receiverId uuid.UUID) ([]EventRow, error) {
	sql := `
	SELECT id, user_id, receiver_id, chat_id, payload FROM messenger.events
	WHERE id > $1 AND receiver_id = $2
	ORDER BY id`
	rows, _ := repository.db.Query(sql, cursor, receiverId)
	events := make([]EventRow, 0)
	for rows.Next() {
		event := EventRow{}
		rows.Scan(
			&event.Id,
			&event.UserId,
			&event.ReceiverId,
			&event.ChatId,
			&event.Payload,
		)
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
