package chat

import (
	"database/sql"

	"github.com/google/uuid"
)

type ChatUsers []uuid.UUID

type ChatRow struct {
	Id    uuid.UUID
	Users []uuid.UUID
}

type ChatRepository interface {
	CreateChat(createChatInfo ChatRow) (ChatRow, error)
	AddUsersToChat(chatId uuid.UUID, users ChatUsers) error
	GetChat(chatId uuid.UUID) (ChatRow, error)
	AreUsersInChat(chatId uuid.UUID, userIds []uuid.UUID) (bool, error)
	GetChatsByUser(userId uuid.UUID) ([]ChatRow, error)
}

type chatRepository struct {
	db *sql.DB
}

// https://blog.thibaut-rousseau.com/blog/sql-transactions-in-go-the-good-way/

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (repository *chatRepository) CreateChat(createChatInfo ChatRow) (ChatRow, error) {
	tx, err := repository.db.Begin()
	if err != nil {
		return ChatRow{}, err
	}
	defer tx.Rollback()

	err = repository.insertChatAndGetId(tx, createChatInfo)
	if err != nil {
		return ChatRow{}, err
	}

	err = repository.addUsersToChatInTx(tx, createChatInfo.Id, createChatInfo.Users)
	if err != nil {
		return ChatRow{}, err
	}
	return createChatInfo, tx.Commit()
}

func (repository *chatRepository) insertChatAndGetId(tx *sql.Tx, chat ChatRow) error {
	sql := `INSERT INTO messenger.chats(шв) VALUES ($1)`
	_, err := tx.Exec(sql, chat.Id)
	return err
}

func (repository *chatRepository) AddUsersToChat(chatId uuid.UUID, users ChatUsers) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = repository.addUsersToChatInTx(tx, chatId, users)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (repository *chatRepository) addUsersToChatInTx(tx *sql.Tx, chatId uuid.UUID, users ChatUsers) error {
	sql := `
	INSERT INTO messenger.chats_users(user_id, chat_id)
	VALUES($1, $2)`
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, userId := range users {
		if _, err := stmt.Exec(userId, chatId); err != nil {
			return err
		}
	}
	return nil
}

func (repository *chatRepository) GetChat(chatId uuid.UUID) (ChatRow, error) {
	getChatSql := `SELECT id FROM messenger.chats WHERE id = $1`
	chat := ChatRow{}
	err := repository.db.QueryRow(getChatSql, chatId).Scan(&chat.Id)
	if err != nil {
		return ChatRow{}, err
	}
	getChatUsersSql := `SELECT user_id FROM messenger.chats_users WHERE chat_id = $1`
	rows, err := repository.db.Query(getChatUsersSql, chatId)
	if err != nil {
		return ChatRow{}, err
	}
	for rows.Next() {
		var userId uuid.UUID
		rows.Scan(&userId)
		chat.Users = append(chat.Users, userId)
	}
	if rows.Err() != nil {
		return ChatRow{}, err
	}
	return chat, nil
}

func (repository *chatRepository) GetChatsByUser(userId uuid.UUID) ([]ChatRow, error) {
	sql := `
	SELECT c.id FROM messenger.chats c
	JOIN messenger.chats_users cu ON c.id = cu.chat_id
	WHERE cu.user_id = $1`
	rows, err := repository.db.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	chats := make([]ChatRow, 0)
	for rows.Next() {
		chatRow := ChatRow{}
		rows.Scan(&chatRow.Id)
		chats = append(chats, chatRow)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return chats, nil
}

func (repository *chatRepository) AreUsersInChat(chatId uuid.UUID, userIds []uuid.UUID) (bool, error) {
	sql := `SELECT COUNT(1) FROM messenger.chats_users WHERE chat_id = $1 AND user_id IN ($2))`
	row := repository.db.QueryRow(sql, chatId, userIds)
	var userInChatCount int
	err := row.Scan(&userInChatCount)
	if err != nil {
		return false, err
	}
	return userInChatCount == len(userIds), nil
}
